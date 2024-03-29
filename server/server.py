import bluetooth
import subprocess
import os

#auto directory setup
if os.path.exists("conin") == False:
    os.mkdir("conin")
if os.path.exists("in") == False:
    os.mkdir("in")
if os.path.exists("out") == False:
    os.mkdir("out")

#connection setup
server_sock = bluetooth.BluetoothSocket(bluetooth.RFCOMM)
port = 5
server_sock.bind(("",port))
server_sock.listen(1)
print("LISTENING ON PORT:", port)

#main loop
while True:
    #wait for client
    print("Waiting for client...")
    client_sock,client_info = server_sock.accept()
    print("ACCEPTED CONNECTION:",client_info)

    #file name info recv
    fname = client_sock.recv(1024).decode('utf-8') #xxx.jpg
    client_sock.send("server: file name recvied")

    #write recv bytes of img
    print("Writing recv bytes...(img.xxx)")
    with open("in/"+fname, 'wb') as f:
        while True:
            data = client_sock.recv(1000)
            if data == b'end':
                break
            f.write(data)
            client_sock.send("server: ok got that chunk")

    #write recv bytes of config.json
    print("Writing recv bytes...(config.json)")
    with open("conin/config.json", 'wb') as f:
        while True:
            data = client_sock.recv(1000)
            if data == b'end':
                break
            f.write(data)
            client_sock.send("server: ok got that chunk")

    #convert image --------------#
    client_sock.send("converting...")
    subprocess.run(['go','run','pixgen.go']) #with this it auto display the subprocess stdout to this stdout
    client_sock.send("finish converting")
    #----------------------------#

    #send file back
    print("Sending converted file back...")
    with open("out/"+fname+".html", 'rb') as f:
        while True:
            data = f.read(1000)
            if data == b'':
                client_sock.send("end")
                break
            client_sock.send(data)
            client_sock.recv(1024)
    #delete sended file
    os.remove("out/"+fname+".html")
    #-----------------------------#

    #finsih this cliend so close connection
    print("Complete this client...")
    client_sock.close()

#close server
server_sock.close()
print("SERVER CLOSED")


