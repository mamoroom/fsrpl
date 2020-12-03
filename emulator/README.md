## Start Emulator

```bash
docker-compose up --build
docker exec -it emulator_firebase_1 bash
firebase login firebase login --no-localhost
firebase init
firebase emulators:start
```
