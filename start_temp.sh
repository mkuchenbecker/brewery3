echo "Starting Temperature Sensors"
sudo modprobe w1-gpio
sudo modprobe w1-therm
echo "ls /sys/bus/w1/devices"
ls /sys/bus/w1/devices
sudo chown root.gpio /dev/gpiomem
sudo chmod g+rw /dev/gpiomem

