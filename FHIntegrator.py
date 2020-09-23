'''
Created on 10 Oct 2019

@author: simon
'''

#!/usr/bin/env python
import pika
import sys, time, json
import subprocess
###
# Fault Handler Integrator
# This is to show how the FH could manage on
# its necessary pub & sub topics with rabbitmq

### Setup of FH Integrator connection
print("## Beginning FH Integrator")
credentials = pika.PlainCredentials('guest', 'password')
connection = pika.BlockingConnection(pika.ConnectionParameters('localhost', 5672, '/', credentials))
channel = connection.channel()
channel.exchange_declare(exchange='topics', exchange_type='topic', durable=True)
key_power = 'Request.Power'
key_monitor = 'Monitor.State'
key_motion = 'Motion.Detected'
key_issue = 'Issue.Notice'
key_failure_component = 'Failure.Component'
key_failure_network   = 'Failure.Network'
key_event = 'Event.FH'
key_email_request = 'Email.Request'
key_email_response = 'Email.Response'
#

# Publishing
result = channel.queue_declare('', exclusive=False, durable=True)
queue_name = result.method.queue
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=key_event)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=key_power)
channel.queue_bind(exchange='topics', queue=queue_name, routing_key=key_email_request)
#
text = '{ "severity": 0, "component": "SYP", "action": null }'
failure = '{ "time":"14:56:00", "type": "Camera", "severity": 1 }'
monitor = '{ "state":false }'
email = '{ "accounts":[ { "role": "ADMIN", "email": "N/A" }] }'
channel.basic_publish(exchange='topics', routing_key=key_issue, body=text)
channel.basic_publish(exchange='topics', routing_key=key_failure_component, body=failure)
channel.basic_publish(exchange='topics', routing_key=key_monitor, body=monitor)
time.sleep(3)
channel.basic_publish(exchange='topics', routing_key=key_email_response, body=email)
#
print("Waiting for Messages")
count = 0
queue_empty = False

def callback(ch, method, properties, body):
    print("Received %r:%r" % (method.routing_key, body))
    print("Count is : ", count)
    time.sleep(0.3)


while not queue_empty:
    method, properties, body = channel.basic_get(queue=queue_name, auto_ack=False)
    if not body is None:
        callback(channel, method, properties, body)
        count = count + 1
       
