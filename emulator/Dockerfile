FROM node:12

RUN yarn global add firebase-tools && \
    firebase setup:emulators:firestore && \
    wget -qO - https://adoptopenjdk.jfrog.io/adoptopenjdk/api/gpg/key/public | apt-key add - && \
    apt update && \
    apt install -y software-properties-common apt-transport-https && \
    add-apt-repository -y https://adoptopenjdk.jfrog.io/adoptopenjdk/deb/ && \
    apt update && \
    apt install -y adoptopenjdk-11-hotspot

COPY conf/firebase.json .

# CMD [ "firebase", "emulators:start"]