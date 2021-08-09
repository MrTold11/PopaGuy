import speech_recognition as sr
import time
import threading
import base64
import requests

r = sr.Recognizer()
r.energy_threshold = 200
counter = 0
AZURE_TOKEN = "token"
IP = "http://192.168.234.16:8080/process"

print("Init OK")

def sendAudio(audio, id):
    req = requests.post(IP, files={"rec" + str(id) + ".d64" : base64.b64encode(audio.get_wav_data())})
#    req = requests.post(IP, files={"rec" : audio.get_wav_data()})
    if req.status_code == 200:
        print("Send OK")
    else:
        print("Send FAIL")

while True:
    with sr.Microphone() as mic:
        audio = r.listen(mic)
    print("Audio recorded")
    t = threading.Thread(target=sendAudio, args = (audio, counter))
    t.start()
    counter += 1
    time.sleep(0.5)
