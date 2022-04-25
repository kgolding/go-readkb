OMG... reading keyboard scancodes is a nightmare!!!

This code will read raw scancodes and decode 'normal' keypresses to a go channel.

A single hardware event can & will send multiple events, these are tied together using `timeval`

The C struct of an event:
```
struct input_event {
	struct timeval time;
	unsigned short type;
	unsigned short code;
	unsigned int value;
};
```

## Linux

### Run the code

1. `cd cmd`
1. `go build && sudo ./cmd /dev/input/event20`

### List all input devices

`cat /proc/bus/input/devices`

```
I: Bus=0003 Vendor=046d Product=c326 Version=0110
N: Name="Logitech USB Keyboard"
P: Phys=usb-0000:00:14.0-11/input0
S: Sysfs=/devices/pci0000:00/0000:00:14.0/usb1/1-11/1-11:1.0/0003:046D:C326.0004/input/input22
U: Uniq=
H: Handlers=sysrq kbd event20 leds 
B: PROP=0
B: EV=120013
B: KEY=1000000000007 ff9f207ac14057ff febeffdfffefffff fffffffffffffffe
B: MSC=10
B: LED=7
```

### View the raw events

Where `event20` comes from the output above.

`sudo cat /dev/input/event20 | hexdump`

### References:
http://www.quadibloc.com/comp/scan.htm
https://www.kernel.org/doc/html/v4.14/input/event-codes.html#input-event-codes
https://elixir.bootlin.com/linux/latest/source/include/uapi/linux/input-event-codes.h#L39
