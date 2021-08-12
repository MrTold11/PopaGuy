from subprocess import call
import speech_recognition as sr
import time, os, random, math, importlib
import threading
import base64
import requests

IP = "http://mrtold.tplinkdns.com/4949" #Server ip
IP_SEND = "process"                #Server path to send audio data
IP_AD = "listad"                   #Server path to get ad list
IP_GETAD = "getad"                 #Server path to get ad file
mute = False                       #Mute parrot?
adMode = False                     #Play ad only (good for dining room / other loud place)
AD_MODE_TO = 5                     #How often to play ads in ad mode (sec)
AD_TIMEOUT = 80                    #How often to play ads (sec)
AD_SYNC_TO = 10                    #How often to sync ads (sec)
AD_DIVIDER = "$$"                  #String that divide ads' names in requested ad list
REC_LIMIT = 15                     #Maximum record duration
MIN_VOLUME = 220                   #Volume floor for record to start
BUFFER_VOICE = 10                  #How many phrases to store?
T_MIN = 5                          #Minimum time to pass before playing phrase (sec)
T_MAX = 210                        #Maximum time to pass before playing phrase (sec) (Divide this by ~4!)
AZURE_TOKEN = "token"              #Token for Azure (todo)
ADS_PATH = "ads"                   #Local path to store ad files
VOICE_PATH = "voice/aud"           #Local path (with part of name) to store phrases
WAV = ".wav"                       #Audio format
OGG = ".ogg"

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
lastPlay = ""

print("Init OK")

# Send recorded audio to server and get new phrase or nothing
def sendAudio(audio, id):
    try:
        if len(audio.frame_data) < 100000:
            print("Short - skip")
            return
        req = requests.post(IP + IP_SEND, files={"rec" + str(id) + ".d64" : (None, base64.b64encode(audio.get_wav_data()))})
        if req.status_code == 200:
            print(" <- Send OK")
            ed = str(req.content)[2:-1]
            if len(ed) < 40:
                print(" -> Nothing")
            else:
                nid = getOldVoice()
                #inrec = ap.AudioProcessing(base64.b64decode(ed))
                #inrec.set_audio_speed(1.2)
                #inrec.save_to_file(VOICE_PATH + str(nid) + WAV)
                with open(VOICE_PATH + str(nid) + WAV, "wb") as ff:
                    ff.write(base64.b64decode(ed))
                print(" -> Phrase recieved")
        else:
            print("XX Send FAIL with code " + str(req.status_code))
    except:
        print("XX An error ocupped during audio sending")
        

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
    print("II Delta: " + str(deltaSpeak) + "; free: " + str(freeSeconds))
    if mute or freeSeconds < 1: return
    if (playAd()):
        return
    if random.randrange(100) < math.sin((deltaSpeak - T_MIN) * MUL_CONST) * 100:
        print("OO Speak after " + str(freeSeconds) + " seconds silence; " + str(deltaSpeak) + " seconds after the last speech")
        freeSeconds = 0
        lastSpeak = time.time()
        fid = "none"
        while True:
            fid = VOICE_PATH + str(random.randrange(BUFFER_VOICE)) + WAV
            if os.path.isfile(fid) and not lastPlay == fid:
                break
            if freeSeconds > 20:
                freeSeconds = 0
                print("Nothing to speak")
                return
            freeSeconds += 1
        freeSeconds = 0
        lastPlay = fid
        call(["aplay", fid, "-D", "hw:0,0"])

# Get ad file from server with name
def getAd(name):
    req = requests.post(IP + IP_GETAD, files={"adName" : (None, str(name))})
    if req.status_code == 200:
        ed = str(req.content)[2:-1]
        if len(ed) < 40:
            print("XX No ad recieved during request")
        else:
            with open(ADS_PATH + "/" + str(name) + OGG, "wb") as f:
                f.write(base64.b64decode(ed))
            print(" -> Ad recieved")
    else:
        print("XX Ad get FAIL with code " + str(req.status_code))

# Get actual ad list and compare it with stored files
def syncAds():
    print(" <- Ad sync...")
    changed = False
    req = requests.post(IP + IP_AD, timeout=2)
    if req.status_code == 200:
        adl = str(req.content)[2:-1].split(AD_DIVIDER)
        for adName in os.listdir(ADS_PATH):
            if adName[:-4] in adl:
                adl.remove(adName[:-4])
            else:
                changed = True
                os.remove(ADS_PATH + "/" + adName)
        if not str(req.content) == "b'NON'":
            for newAd in adl:
                changed = True
                getAd(newAd)
        if changed:
            listAds()
    else:
        print("XX Ad sync FAIL with code " + str(req.status_code))

def listAds():
    global adStat
    msa = 100000
    for oad in adStat:
        msa = min(msa, adStat[oad])
    if msa == 100000:
        msa = 0
    oldStat = adStat
    adStat = {}
    for adName in os.listdir(ADS_PATH):
        adStat[adName] = oldStat.get(adName, msa)

# Play ad lol
def playAd():
    global lastAd, lastAdFr
    if int(time.time() - lastAd) < AD_MODE_TO:
        return False
    if not adMode:
        if mute or int(time.time()) - lastAd < AD_TIMEOUT:
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
        print("OO Playing ad...")
        call(["mplayer", "-ao", "alsa:noblock:device=hw=0.0", "-really-quiet", ADS_PATH + "/" + str(nad)])
        adStat[nad] = adStat[nad] + 1
        lastAd = time.time()
        return True
    return False

def adSyncLoop():
    while True:
        try:
            syncAds()
        except:
            print("XX Error during ad sync")
        time.sleep(AD_SYNC_TO)

listAds()
adSyncThread = threading.Thread(target=adSyncLoop)
adSyncThread.start()

while True:
    if adMode:
        playAd()
        time.sleep(1)
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
