import azure.cognitiveservices.speech as speechsdk
import base64
import wave
from model import predict
from azure.cognitiveservices.speech import AudioDataStream, SpeechConfig, SpeechSynthesizer
from azure.cognitiveservices.speech.audio import AudioOutputConfig
speech_config = speechsdk.SpeechConfig(subscription="29103316047b4c59a3edbadb3c899c81", endpoint="https://northeurope.api.cognitive.microsoft.com/sts/v1.0/issuetoken")
speech_config.speech_recognition_language = 'ru-RU'
speech_config.speech_synthesis_language = 'ru-RU'
#speech_config.speech_synthesis_voice_name = 'ru-RU-DmitryNeural'

inf = open("sound.b64")
for i in range(3): inf.readline()
byt = base64.b64decode(inf.readline())
with open("NEWSOUND.wav", 'wb') as f:
    f.write(byt)



numbers = {'1':'один',
           "2": "два",
           "3": "три",
           '4': 'четыре',
           "5": "пять",
           "6": "шесть",
           '7': 'семь',
           "8": "восемь",
           "9": "девять",
           "10": "десять"}



def from_file(file_name):
    print('from file')
    audio_input = speechsdk.AudioConfig(filename=file_name)
    speech_recognizer = speechsdk.SpeechRecognizer(speech_config=speech_config, audio_config = audio_input)
    result = speech_recognizer.recognize_once_async().get()
    return result.text


def gen_voice(text, output):
    result = []
    for i in range(len(text)):
        s = """<speak version="1.0" xmlns="https://www.w3.org/2001/10/synthesis" xml:lang="ru-RU">
                    <voice name="ru-RU-Irina">
                        <prosody volume="+40%" pitch="1000Hz" contour="(60%,-60%) (100%,+80%)">
                            """+ text[i] +"""
                        </prosody>
                    </voice>
            </speak>"""
        with open('style.xml', 'w') as f:
            f.write(s)
#        audio_config = AudioOutputConfig(filename="sentence" + str(i) + '.wav')
        synthesizer = SpeechSynthesizer(speech_config=speech_config, audio_config=None)
        ssml_string = open('style.xml', 'r').read()
        result = synthesizer.speak_ssml_async(ssml_string).get()
        stream = AudioDataStream(result)
#        synthesizer.speak_text_async(text[i])
        stream.save_to_wav_file("sentence" + str(i) + ".wav")

sentences = from_file("NEWSOUND.wav").replace('?', '.').replace('!', '.').split('.')
text = []
for i in sentences:
    if i == '': continue
    buf = i.lower().replace('!', '').replace('.', '').replace('?','').replace(',', '')
    for j in buf.split():
        if j in numbers:
            buf = buf.replace(j, numbers[j])
    t = predict(buf)
    if -20 < t < -10:
        print(t)
        text.append(i)
print('len: ' + str(len(text)))
with open("number.txt", 'w') as f:
    f.write(str(len(text)))
gen_voice(text, "/home/mrtold/popaguy/PopaGuy/")
