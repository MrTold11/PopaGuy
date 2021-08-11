from subprocess import call
import speech_recognition as sr
#import AudioProcessing as ap
import time, os, random, math, importlib
import threading
import base64
import requests

IP = "http://192.168.217.88:8080/" #Server ip
IP_SEND = "process"                #Server path to send audio data
IP_AD = "listad"                   #Server path to get ad list
IP_GETAD = "getad/"                #Server path to get ad file (IP_GETAD + name)
mute = False                       #Mute parrot?
adMode = False                     #Play ad only (good for dining room / other loud place)
AD_TIMEOUT = 60                    #How often to play ads (sec)
AD_DIVIDER = "$$"                  #String that divide ads' names in requested ad list
PROCESS_NOISE = False               #Remove noice and increase volume?
REC_LIMIT = 15                     #Maximum record duration
MIN_VOLUME = 250                   #Volume floor for record to start
BUFFER_VOICE = 10                  #How many phrases to store?
T_MIN = 10                         #Minimum time to pass before playing phrase (sec)
T_MAX = 420                        #Maximum time to pass before playing phrase (sec) (Divide this by ~4!)
AZURE_TOKEN = "token"              #Token for Azure (todo)
ADS_PATH = "ads"                   #Local path to store ad files
VOICE_PATH = "voice/aud"           #Local path (with part of name) to store phrases
WAV = ".wav"                       #Audio format

r = sr.Recognizer()
r.energy_threshold = MIN_VOLUME
r.dynamic_energy_threshold = False
MUL_CONST = 1.570796326794896 / (T_MAX - T_MIN)
counter = 0
freeSeconds = 0
lastSpeak = time.time()
lastAd = time.time()
lastAdFr = 0
random.seed(lastSpeak)
adStat = {}

#importlib.import_module("AudioProcessing")

print("Init OK")

# Send recorded audio to server and get new phrase or nothing
def sendAudio(audio, id):
    if len(audio.frame_data) < 200000:
        print("Short - skip")
        return
    #recf = ap.AudioProcessing(audio.get_wav_data())
    #if PROCESS_NOISE:
    #    recf.processNoise(2, 0.7)
    #recf.save_to_file("out/rec" + str(id) + WAV)
    with open("out/rec" + str(id) + WAV, "rb") as rf:
        req = requests.post(IP + IP_SEND, files={"rec" + str(id) + ".d64" : (None, base64.b64encode(audio.get_wav_data()))})
    if req.status_code == 200:
        print("Send OK")
        ed = str(req.content)[2:-1]
        if len(ed) < 40:
            print("Nothing")
        else:
            nid = getOldVoice()
            #inrec = ap.AudioProcessing(base64.b64decode(ed))
            #inrec.set_audio_speed(1.2)
            #inrec.save_to_file(VOICE_PATH + str(nid) + WAV)
            with open(VOICE_PATH + str(nid) + WAV, "wb") as ff:
                ff.write(base64.b64decode(ed))
            print("Phrase recieved")
#            call(["aplay", VOICE_PATH + str(nid) + WAV, "-D", "hw:0,0"])
#            print("Playing...")
    else:
        print("Send FAIL")

# Return id of the oldest file stored in VOICE_PATH
def getOldVoice():
    of = 10
    ot = int(time.time())
    for id in range(0, BUFFER_VOICE):
        fp = VOICE_PATH + str(id) + WAV
        if not os.path.isfile(fp):
            return id
        lm = int(os.path.getmtime(fp))
        if lm < ot:
            of = id
            ot = lm
    return of

# Try playing any sound (weird probability)
def tryPlay():
    global freeSeconds, lastSpeak
    deltaSpeak = int(time.time() - lastSpeak)
    print("Delta: " + str(deltaSpeak) + "; free: " + str(freeSeconds))
    if mute or freeSeconds < 1: return
    if (playAd()):
        return
    if random.randrange(100) < math.sin((deltaSpeak - T_MIN) * MUL_CONST) * 100:
        print("Speak after " + str(freeSeconds) + " seconds silence; " + str(deltaSpeak) + " seconds after the last speech")
        freeSeconds = 0
        lastSpeak = time.time()
        fid = "none"
        while True:
            fid = VOICE_PATH + str(random.randrange(BUFFER_VOICE)) + WAV
            if os.path.isfile(fid):
                break
            if freeSeconds > 20:
                freeSeconds = 0
                print("Nothing to speak")
                return
            freeSeconds += 1
        freeSeconds = 0
        call(["aplay", fid, "-D", "hw:0,0"])

# Get ad file from server with name
def getAd(name):
    req = requests.post(IP + IP_GETAD + name)
    if req.status_code == 200:
        adl = str(req.content)[2:-1]
        if len(ed) < 40:
            print("No ad recieved during request")
        else:
            with open(ADS_PATH + "/" + str(name) + WAV, "wb") as f:
                f.write(base64.b64decode(ed))
            print("Ad recieved")
    else:
        print("Ad get FAIL")

# Get actual ad list and compare it with stored files
def syncAds():
    print("Ad sync...")
    changed = False
    req = requests.post(IP + IP_AD)
    if req.status_code == 200:
        adl = str(req.content)[2:-1].split(AD_DIVIDER)
        for adName in os.listdir(ADS_PATH):
            if adName[:-4] in adl:
                adl.remove(adName[:-4])
            else:
                changed = True
                os.remove(ADS_PATH + "/" + adName)
        
        for newAd in adl:
            changed = True
            getAd(newAd)
        if changed:
            listAds()
    else:
        print("Ad sync FAIL")

def listAds():
    global adStat
    msa = 100000
    for oad in adStat:
        msa = min(msa, adStat[oad])
    oldStat = adStat
    adStat = {}
    for adName in os.listdir(ADS_PATH):
        adStat[adName] = oldStat.get(adName, msa)

# Play ad lol
def playAd():
    global lastAd, lastAdFr
    if int(time.time() - lastAd) < AD_TIMEOUT:
        return False
    if mute and not adMode:
        return False
    minAd = 10000
    nad = None
    for adName in os.listdir(ADS_PATH):
        ca = adStat[adName]
        if minAd > ca:
            nad = adName
            minAd = ca
    if nad == None:
        print("No ad to play")
    else:
        print("Playing ad...")
        call(["aplay", ADS_PATH + "/" + str(nad), "-D", "hw:0,0"])
        adStat[nad] = adStat[nad] + 1
        lastAd = time.time()
        return True
    return False

while True:
    if adMode:
        playAd()
        sleep(1)
        continue
    try:
        with sr.Microphone() as mic:
            audio = r.listen(mic, 1, REC_LIMIT)
        if random.randrange(100) < 40:
            tryPlay()
    except sr.WaitTimeoutError:
        freeSeconds += 1
        tryPlay()
        continue
    freeSeconds = 0
    print("Audio recorded")
    t = threading.Thread(target=sendAudio, args = (audio, counter))
    t.start()
    counter += 1
    time.sleep(0.5)
