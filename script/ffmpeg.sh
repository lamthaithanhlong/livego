#!/bin/bash

while true; do
    # Start the FFmpeg command
    ffmpeg -i rtmp://localhost/live/livestream -c copy -f flv rtmp://167.88.168.20:1935/live/test 2>ffmpeg_error.log &

    FFMPEG_PID=$!
    echo "Stream started with PID $FFMPEG_PID"

    # Sleep for 180 seconds (3 minutes)
    sleep 180

    # Check if FFmpeg process is running
    if ps -p $FFMPEG_PID > /dev/null
    then
        echo "180 seconds passed, restarting stream..."
        kill $FFMPEG_PID
    else
        echo "FFmpeg exited before 180 seconds."
    fi

    # Check for errors in FFmpeg output
    if grep -q 'error\|freeze' ffmpeg_error.log; then
        echo "Error or freeze detected in FFmpeg output, restarting stream..."
    fi

    # Short delay before restarting
    sleep 1
done
