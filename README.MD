# gobot-fat

## Build

```bash
GOOS=linux GOARCH=arm GOARM=6 go build

GOOS=linux GOARCH=arm GOARM=6 CGO_ENABLED=1 CC=arm-linux-gnueabihf-gcc go build
```

## RUN

```bash
docker run -d --name dfp -ti -p 4040:4040 -v /opt/dfp/data:/opt/dfp/data -v /opt/dfp/config:/opt/dfp/config disaster37/dfp:latest 
```

## Handle Technical Filter Pond Robot (TFP robot)

### Login
```bash
curl -XPOST -H "Content-Type: application/json" http://localhost:4040/token-auth  -d '{"username": "gobot", "password": "gobot"}'
```

### Start pond pump
```bash
curl -XPOST -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG1pbiI6dHJ1ZSwiZXhwIjoxNTkzOTU3MTk2LCJuYW1lIjoiZ29ib3QifQ.dnblWfYBS3tQ2iokTPmEJpYzCKW8UB-FdXM-KQbMWvI" http://localhost:4040/api/tfps/action/start_pond_pump
```

### Stop pond pump
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/stop_pond_pump
```

### Start pond pump with UVC
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/start_pond_pump_with_uvc
```

### Start waterfall pump
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/start_waterfall_pump
```

### Stop waterfall pump
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/stop_waterfall_pump
```

### Start UVC1
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/start_uvc1
```

### Stop UVC1
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/stop_uvc1
```

### Start UVC2
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/start_uvc2
```

### Stop UVC2
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/stop_uvc2
```

### Start pond bubble
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/start_pond_bubble
```

### Stop pond bubble
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/stop_pond_bubble
```

### Start filter bubble
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/start_filter_bubble
```

### Stop filter bubble
```bash
curl -XPOST -u gobot:gobot http://localhost:4040/api/tfp/action/stop_filter_bubble
```

### Get UVC1 blister status
```bash
curl -XGET -u gobot:gobot http://localhost:4040/api/tfp/uvc/uvc1_blister_status
```

### Get UVC2 blister status
```bash
curl -XGET -u gobot:gobot http://localhost:4040/api/tfp/uvc/uvc2_blister_status
```

### Put UVC1 blister change
```bash
curl -XPUT -u gobot:gobot http://localhost:4040/api/tfp/uvc/uvc1_blister_new
```

### Put UVC2 blister change
```bash
curl -XPUT -u gobot:gobot http://localhost:4040/api/tfp/uvc/uvc2_blister_new
```

