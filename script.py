import subprocess
import time
from PIL import Image
import numpy as np
from ping import send_ipv6_ping

def read_image(path):
    # read the png file from the path and get a array of pixels
    image = Image.open(path)
    if image.mode != 'RGB':
        image = image.convert('RGB')
    pixel_array = np.array(image)

    
    return image.size, pixel_array





PRE = '2001:610:1908:a000:'

import subprocess


def send_ping6(ip):
    try:
        process = subprocess.Popen(['ping6', '-c', '1', ip], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    except FileNotFoundError:
        print('Ping6 not found')

def draw(x, y, r, g, b):
    x = int(x)
    y = int(y)
    r = int(r)
    g = int(g)
    b = int(b)



    a = 255
    #<PRE>:<X>:<Y>:<B><G>:<R><A>
    x = f"{x:04x}"
    y = f"{y:04x}"

    b = f"{b:02x}"
    g = f"{g:02x}"
    r = f"{r:02x}"
    a = f"{a:02x}"
    ip = f"{PRE}{x}:{y}:{b}{g}:{r}{a}"
    send_ipv6_ping(ip, 1)
    


screen_w = 1920
screen_h = 1080








queue = []

def draw_image(path, x2, y2):
    global queue
    size, pixels = read_image(path)
    for y in range(size[1]):
        for x in range(size[0]):
            r, g, b = pixels[y][x]
            queue.append((x2+x, y2+y, r, g, b))




startX = int(screen_w/2)
startY = int(screen_h/2)

draw_image('vold.jpg', 10, 10)




batch_size = 100

total = len(queue)
for i in range(0, len(queue), batch_size):
    batch = queue[i:i+batch_size]
    for x, y, r, g, b in batch:
        draw(x, y, r, g, b)
    
    if i % 1000 == 0:
        progress = i/total * 100 
        print(f"Sent {i+batch_size}/{total} packets ({progress:.2f}%)")
    time.sleep(0.01)


time.sleep(4)
