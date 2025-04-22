#include <fcntl.h>
#include <linux/spi/spidev.h>
#include <sys/ioctl.h>
#include <unistd.h>
#include <cstring>
#include <cstdio>
#include <vector>
#include <bcm2835.h>

#define SPI_DEVICE_0 "/dev/spidev0.0"
#define SPI_DEVICE_1 "/dev/spidev0.1"
#define CS0 RPI_GPIO_P1_24 // GPIO 8
#define CS1 RPI_GPIO_P1_26 // GPIO 7

const uint8_t CMD_SOFT_RESET = 0xCC;
const uint8_t CMD_GLOBAL_BRIGHTNESS = 0xFF;
const uint8_t CMD_COM_PIN_CTRL = 0x41;
const uint8_t CMD_ROW_PIN_CTRL = 0x42;
const uint8_t CMD_WRITE_DISPLAY = 0x80;
const uint8_t CMD_SYSTEM_CTRL = 0x35;
const uint8_t CMD_SCROLL_CTRL = 0x20;

const int COLS = 17;
const int ROWS = 7;
const int NUM_PIXELS = COLS * ROWS;
const int BUF_SIZE = 28 * 8 * 2;

int spi_fd0, spi_fd1;
std::vector<uint8_t> disp(NUM_PIXELS * 3, 0);
std::vector<uint8_t> buf(BUF_SIZE, 0);

std::vector<std::vector<int>> lut = {
    {139, 138, 137}, {223, 222, 221}, {167, 166, 165}, {195, 194, 193}, {111, 110, 109}, {55, 54, 53}, {83, 82, 81},
    {136, 135, 134}, {220, 219, 218}, {164, 163, 162}, {192, 191, 190}, {108, 107, 106}, {52, 51, 50}, {80, 79, 78},
    {113, 115, 114}, {197, 199, 198}, {141, 143, 142}, {169, 171, 170}, {85, 87, 86}, {29, 31, 30}, {57, 59, 58},
    {116, 118, 117}, {200, 202, 201}, {144, 146, 145}, {172, 174, 173}, {88, 90, 89}, {32, 34, 33}, {60, 62, 61},
    {119, 121, 120}, {203, 205, 204}, {147, 149, 148}, {175, 177, 176}, {91, 93, 92}, {35, 37, 36}, {63, 65, 64},
    {122, 124, 123}, {206, 208, 207}, {150, 152, 151}, {178, 180, 179}, {94, 96, 95}, {38, 40, 39}, {66, 68, 67},
    {125, 127, 126}, {209, 211, 210}, {153, 155, 154}, {181, 183, 182}, {97, 99, 98}, {41, 43, 42}, {69, 71, 70},
    {128, 130, 129}, {212, 214, 213}, {156, 158, 157}, {184, 186, 185}, {100, 102, 101}, {44, 46, 45}, {72, 74, 73},
    {131, 133, 132}, {215, 217, 216}, {159, 161, 160}, {187, 189, 188}, {103, 105, 104}, {47, 49, 48}, {75, 77, 76},
    {363, 362, 361}, {447, 446, 445}, {391, 390, 389}, {419, 418, 417}, {335, 334, 333}, {279, 278, 277},
    {307, 306, 305}, {360, 359, 358}, {444, 443, 442}, {388, 387, 386}, {416, 415, 414}, {332, 331, 330},
    {276, 275, 274}, {304, 303, 302}, {337, 339, 338}, {421, 423, 422}, {365, 367, 366}, {393, 395, 394},
    {309, 311, 310}, {253, 255, 254}, {281, 283, 282}, {340, 342, 341}, {424, 426, 425}, {368, 370, 369},
    {396, 398, 397}, {312, 314, 313}, {256, 258, 257}, {284, 286, 285}, {343, 345, 344}, {427, 429, 428},
    {371, 373, 372}, {399, 401, 400}, {315, 317, 316}, {259, 261, 260}, {287, 289, 288}, {346, 348, 347},
    {430, 432, 431}, {374, 376, 375}, {402, 404, 403}, {318, 320, 319}, {262, 264, 263}, {290, 292, 291},
    {349, 351, 350}, {433, 435, 434}, {377, 379, 378}, {405, 407, 406}, {321, 323, 322}, {265, 267, 266},
    {293, 295, 294}, {352, 354, 353}, {436, 438, 437}, {380, 382, 381}, {408, 410, 409}, {324, 326, 325},
    {268, 270, 269}, {296, 298, 297}
};

void toggle_cs(uint8_t cs_pin, bool active) {
    bcm2835_gpio_write(cs_pin, active ? LOW : HIGH);
}

bool spi_transfer(int fd, uint8_t cs_pin, const std::vector<uint8_t>& data) {
    toggle_cs(cs_pin, true);
    if (write(fd, data.data(), data.size()) != (ssize_t)data.size()) {
        perror("SPI write failed");
        toggle_cs(cs_pin, false);
        return false;
    }
    toggle_cs(cs_pin, false);
    return true;
}

void send_cmds(int fd, uint8_t cs_pin, const std::vector<std::vector<uint8_t>>& cmds) {
    for (auto& cmd : cmds)
        spi_transfer(fd, cs_pin, cmd);
}

void setup_matrix(int fd, uint8_t cs_pin, int offset) {
    send_cmds(fd, cs_pin, {
        {CMD_SOFT_RESET},
        {CMD_GLOBAL_BRIGHTNESS, 0x01},
        {CMD_SCROLL_CTRL, 0x00},
        {CMD_SYSTEM_CTRL, 0x00},
        {CMD_WRITE_DISPLAY, 0x00} // buffer later
    });

    std::vector<uint8_t> data = {CMD_WRITE_DISPLAY, 0x00};
    data.insert(data.end(), buf.begin() + offset, buf.begin() + offset + (28 * 8));
    spi_transfer(fd, cs_pin, data);

    send_cmds(fd, cs_pin, {
        {CMD_COM_PIN_CTRL, 0xff},
        {CMD_ROW_PIN_CTRL, 0xff, 0xff, 0xff, 0xff},
        {CMD_SYSTEM_CTRL, 0x03}
    });
}

void set_all(uint8_t r, uint8_t g, uint8_t b) {
    for (int i = 0; i < NUM_PIXELS; i++) {
        disp[i * 3 + 0] = r >> 2;
        disp[i * 3 + 1] = g >> 2;
        disp[i * 3 + 2] = b >> 2;
    }
}

void show() {
    for (int i = 0; i < NUM_PIXELS; i++) {
        int ir = lut[i][0], ig = lut[i][1], ib = lut[i][2];
        buf[ir] = disp[i * 3];
        buf[ig] = disp[i * 3 + 1];
        buf[ib] = disp[i * 3 + 2];
    }

    std::vector<uint8_t> left_data = {CMD_WRITE_DISPLAY, 0x00};
    left_data.insert(left_data.end(), buf.begin(), buf.begin() + (28 * 8));
    spi_transfer(spi_fd0, CS0, left_data);

    std::vector<uint8_t> right_data = {CMD_WRITE_DISPLAY, 0x00};
    right_data.insert(right_data.end(), buf.begin() + (28 * 8), buf.end());
    spi_transfer(spi_fd1, CS1, right_data);
}

int main() {
    if (!bcm2835_init()) {
        fprintf(stderr, "bcm2835_init failed\n");
        return 1;
    }

    bcm2835_gpio_fsel(CS0, BCM2835_GPIO_FSEL_OUTP);
    bcm2835_gpio_fsel(CS1, BCM2835_GPIO_FSEL_OUTP);
    bcm2835_gpio_write(CS0, HIGH);
    bcm2835_gpio_write(CS1, HIGH);

    spi_fd0 = open(SPI_DEVICE_0, O_WRONLY);
    spi_fd1 = open(SPI_DEVICE_1, O_WRONLY);
    if (spi_fd0 < 0 || spi_fd1 < 0) {
        perror("SPI device open failed");
        return 1;
    }

    int speed = 600000;
    ioctl(spi_fd0, SPI_IOC_WR_MAX_SPEED_HZ, &speed);
    ioctl(spi_fd1, SPI_IOC_WR_MAX_SPEED_HZ, &speed);

    setup_matrix(spi_fd0, CS0, 0);
    setup_matrix(spi_fd1, CS1, 28 * 8);

    while(true) {
        set_all(255, 0, 0); // red
        show();
        usleep(1000000);
    
        set_all(0, 255, 0); // green
        show();
        usleep(1000000);
    
        set_all(0, 0, 255); // blue
        show();
        usleep(1000000);
    
        set_all(0, 0, 0); // off
        show();
    }

    close(spi_fd0);
    close(spi_fd1);
    bcm2835_close();

    return 0;
}
