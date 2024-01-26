import serial
import time

BAUDRATE = 19200
PORTNAME = '/dev/ttyUSB0'

# send command to get raw status string from LM-4
def get_status():
    with serial.Serial(port=PORTNAME, baudrate=BAUDRATE, timeout=None) as COM_port:

        # command 'set wavelength'
        COM_port.write(b'5201\r')
        time.sleep(0.1)
        
        responce = COM_port.read_all()
        return responce

# get current wavelength from raw status string
def convert_status_to_wavelength(responce):
    # Достаем Hex коды длины волны из последовательности 
    # (разбиваем строку на массив из 5 кодов)
    data = str(responce)
    try:
        data = data.split(':')[1][:-5]
    except:
        data = None
    
    if not data:
        return 0
    data = data.split()

    # длина волны содержится в 1 и 2м Hex коде с обратной нотацией
    # Соединяем их в одно Hex число
    wavelength_str = data[2] + data[1]
    # переводим строку с хекс числом в массив байт 
    wavelength = bytearray.fromhex(wavelength_str)
    # из массива байт достаем int число
    wavelength = int.from_bytes(wavelength, byteorder='big')
    return wavelength

# send commant to go to wavelength
def go_to(wavelength):
    with serial.Serial(port=PORTNAME, baudrate=BAUDRATE) as COM_port:
        wl_hex = hex(wavelength)[2:]
        wl_bytes = (wl_hex[2:4] + wl_hex[0:2]).upper()
        command = f'4C{wl_bytes}03\r'

        # command 'set wavelength'
        COM_port.write(bytes(command, 'UTF-8'))

        res = COM_port.read_all()

        print(res)



if __name__ == "__main__":
    while True: 
        # read wavelength from keyboard
        wl = input("select wavelength: ")

        # close if '\end'
        if wl == '\end':
            break
        wl = wl.strip()

        # get current wl if '\get'
        if wl == '\get':
            res = get_status()
            res = convert_status_to_wavelength(res)
            print(res)
            continue
        
        # make some value validation
        if len(wl) < 4:
            print("write 4 digit wavelength")
            continue

        # Angstrom from nanometers
        wl = float(wl) * 10

        # go to set wl
        go_to(int(wl))


    print('End of programm')


