import azure.cognitiveservices.speech as speechsdk
import base64
import wave
from azure.cognitiveservices.speech import AudioDataStream, SpeechConfig, SpeechSynthesizer
from azure.cognitiveservices.speech.audio import AudioOutputConfig
speech_config = speechsdk.SpeechConfig(subscription="29103316047b4c59a3edbadb3c899c81", endpoint="https://northeurope.api.cognitive.microsoft.com/sts/v1.0/issuetoken")
speech_config.speech_recognition_language = 'ru-RU'
speech_config.speech_synthesis_language = 'ru-RU'

inf = open("E:/AD/b64/test.b64")
for i in range(3): inf.readline()
byt = base64.b64decode(inf.readline())
#print(wave.open("E:/AD/dict/Voice2.wav").getframerate())
with open("E:/AD/dict/rec.wav", 'wb') as f:
    f.write(byt)




def from_file(file_name):

    audio_input = speechsdk.AudioConfig(filename=file_name)
    speech_recognizer = speechsdk.SpeechRecognizer(speech_config=speech_config, audio_config = audio_input)
    result = speech_recognizer.recognize_once_async().get()
    return result.text


def gen_voice(text, output):
    result = []
    for i in range(len(text)-1):

        audio_config = AudioOutputConfig(filename=output + "sentence" + str(i) + '.wav')
        synthesizer = SpeechSynthesizer(speech_config=speech_config, audio_config=audio_config)
        synthesizer.speak_text_async(text[i])

sentences = from_file("E:/AD/dict/rec.wav").replace('?', '.').replace('!', '.').split('.')
#a =[]
#for x in sentences:
#    a += x.split('?')
#for x in sentences:
#    a += x.split('!')
#sentences = a
with open("E:/AD/num.txt", 'w') as f:
    f.write(str(len(sentences)-1))
gen_voice(sentences, "E:/AD/result/")
