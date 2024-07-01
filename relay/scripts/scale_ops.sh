#!/bin/bash

operation=$1
export TEST=$1
export OPS=500
ops=500
scale=4
cutoff=80
if [ $1 == "signing" ]; then
    scale=26
    cutoff=40
fi
total_ops=$(echo "$scale * $ops" | bc)
total_scale=0
echo "Running $1 operations in scale.."
for ((i=1; i<=$scale; i++)); do
    output=$(timeout $cutoff ./relay/bin/client_func  2>&1)
    success=$(echo "$output" | grep -o '\.' | wc -l)
    total_scale=$(echo "$total_scale + $success" | bc)
    echo .
    sleep 1
done

echo "Total Throughput scale : $total_scale"