
<h2 align="center">stream2videointelligence</h2>
<p align="center">A simple and effective tool to access the Google Cloud Platform's Video Intelligence API and get near real-time insights</p>
<p align="center">
    <a href="#about">About</a> •
    <a href="#features">Features</a> •
    <a href="#build">Build</a> •
    <a href="#configuration">Configuration</a> •
    <a href="#usage">Usage</a> •
    <a href="#license">License</a>
</p>

## About

Cloud Video Intelligence API Streaming is a simple and efficient tool for streaming to the Google Cloud Platform Video Intelligence API.

## Features

- ✅ A **single** statically compiled **binary** for each OS/architecture
- ✅ Can **retrieve the data** to be streamed **directly via a pipe**
- ✅ Can **retrieve the data** to be streamed **from a file**
- ✅ Can **save annotations** in a **Cloud Storage bucket**
- ✅ Can **save annotations** in **real-time** to a **local file**
- ✅ Can **pass annotation** data in **real-time** to another program through its **stdout output**

## Build
First check that you have **Golang** installed on your machine.  
Then, **run**:  
```bash
make 
```
Quite simply!

## Configuration
This program has a number of options, as follows:

| Options    | Description                                          | Mandatory |
| ---------- | ---------------------------------------------------- | --------- |
| `-creds`   | Service account JSON key file path                   | Yes       |
| `-source`  | Path of a file used as a source instead of a pipe    | No        |
| `-feature` | API Cloud Video Intelligence streaming feature       | Yes       |
| `-gcs`     | GCS URI to store all annotation results              | No        |
| `-stdout`  | Print in stdout results from the API                 | No        |
| `-export`  | Export the annotation results from the API to a file | No        |

**Note**:
- You must choose to enable at least one option between `-gcs`, `-stdout` or `-export`
- If `-source` is not set, the data must be piped from another program
- The different features supported by the Cloud Video Intelligence API in streaming mode are `STREAMING_OBJECT_TRACKING`, `STREAMING_LABEL_DETECTION`, `STREAMING_EXPLICIT_CONTENT_DETECTION`, `STREAMING_SHOT_CHANGE_DETECTION`

## Usage

Example of use with a pipe:

```
gst-launch-1.0 -q -v v4l2src device=/dev/video0 ! video/x-raw,framerate=30/1,width=640,height=480 ! videoconvert ! videoscale ! videorate ! x264enc tune=zerolatency speed-preset=6 ! flvmux ! queue ! fdsink | ./stream2videointelligence -creds ./credientials.json -feature STREAMING_OBJECT_TRACKING -stdout
```

Example of use with a video file:
```
./stream2videointelligence -creds ./credentials.json -source video.flv -feature STREAMING_OBJECT_TRACKING -stdout
```

## License

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see http://www.gnu.org/licenses/.
