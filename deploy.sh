docker pull 192.168.1.30:5000/common/xftts:latest
docker rm -f tts-service
docker run -d -p 6061:20000 --name=tts-service --restart=always \
-v ${PWD}/xftts/include:/app/xf/include \
-v ${PWD}/xftts/libs:/app/xf/libs \
-v ${PWD}/xftts/bin/msc:/app/xf/msc \
192.168.1.30:5000/common/xftts:latest -appid=5ff5193f