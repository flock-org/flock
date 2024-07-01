start=$(date +%s%N)
killall wgproxy
wgproxy 8000 54.177.190.187:9000 &
sudo ip link del wireguard1
if [ $1 == "azure" ]; then
    sudo tunnel-benchmarking --tunnel-type=wireguard \
    --datapath=linux --host-name=host1 \
    --remote-hosts=host2:127.0.0.1 \
    --wireguard-public-key=ARq3ziAZVWj5IQ202TskQEsl3GQQrQ7NnJAKOv2F5kE= \
    --wireguard-private-key=uCmjq4myg7GGCZP6Shu7xXuyVzyeyedg/VhZrVJtck4= 
else   
    sudo ip link del wireguard2
    sudo tunnel-benchmarking --tunnel-type=wireguard \
    --datapath=linux --host-name=host2 \
    --remote-hosts=host1:127.0.0.1 \
    --wireguard-public-key=/tcWr8BES3jzSbE2vGxH+PAWqEawE2/2tWe7+LUVHGU= \
    --wireguard-private-key=MOp5jPFweMpLPNSgkp4CLMWzJO2Yh7ASlIQpPAQxwXA= 
fi
end=$(date +%s%N)
time_diff=$(( ($end - $start) / 1000000 ))
echo "Setup Latency: $time_diff ms"
