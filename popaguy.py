from subprocess import call
import speech_recognition as sr
import time, os, random, math
import threading
import base64
import requests

BUFFER_VOICE = 10
T_MIN = 10
T_MAX = 420 #Divide this by ~4
MUL_CONST = 1.570796326794896 / (T_MAX - T_MIN)
AZURE_TOKEN = "token"
IP = "http://192.168.234.16:8080/process"
#IP = "http://192.168.234.39:4455/"
VOICE_PATH = "voice/aud"
VOICE_FORM = ".wav"

r = sr.Recognizer()
r.energy_threshold = 200
counter = 0
freeSeconds = 0
lastSpeak = time.time()
mute = False
random.seed(lastSpeak)

print("Init OK")

def sendAudio(audio, id):
    with open("out/rec" + str(id) + ".wav", "wb") as f:
        f.write(audio.get_wav_data())
    req = requests.post(IP, files={"rec" + str(id) + ".d64" : (None, base64.b64encode(audio.get_wav_data()))})
#    req = requests.post(IP, files={"rec" : audio.get_wav_data()})
    if req.status_code == 200:
        print("Send OK")
        ed = str(req.content)[2:-1]
        if len(ed) < 20:
            print("Nothing")
        else:
            nid = getOldVoice()
            with open(VOICE_PATH + str(nid) + VOICE_FORM, "wb") as f:
                f.write(base64.b64decode(ed))
            print("Phrase recieved")
            call(["aplay", VOICE_PATH + str(nid) + VOICE_FORM, "-D", "hw:0,0"])
            print("Playing...")
    else:
        print("Send FAIL")

def getOldVoice():
    of = 10
    ot = int(time.time())
    for id in range(0, BUFFER_VOICE):
        fp = VOICE_PATH + str(id) + VOICE_FORM
        if not os.path.isfile(fp):
            return id
        lm = int(os.path.getmtime(fp))
        if lm < ot:
            of = id
            ot = lm
    return of

def tryPlay():
    global freeSeconds, lastSpeak
    deltaSpeak = int(time.time() - lastSpeak)
    if mute or freeSeconds < 5: return
    if random.randrange(100) < math.sin((deltaSpeak - T_MIN) * MUL_CONST) * 100:
        print("Speak after " + str(freeSeconds) + " seconds silence; " + str(deltaSpeak) + " seconds after the last speech")
        freeSeconds = 0
        lastSpeak = time.time()
        fid = "none"
        while True:
            fid = VOICE_PATH + str(random.randrange(BUFFER_VOICE)) + VOICE_FORM
            if os.path.isfile(fid):
                break
            if freeSeconds > 20:
                freeSeconds = 0
                print("Nothing to speak")
                return
            freeSeconds += 1
        freeSeconds = 0
        call(["aplay", fid, "-D", "hw:0,0"])

while True:
    try:
        with sr.Microphone() as mic:
            audio = r.listen(mic, 2)
    except WaitTimeoutError:
        freeSeconds += 2
        tryPlay()
        continue
    freeSeconds = 0
    print("Audio recorded")
    t = threading.Thread(target=sendAudio, args = (audio, counter))
    t.start()
    counter += 1
    time.sleep(0.5)
