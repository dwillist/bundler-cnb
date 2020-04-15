# minimal ruby only server for integration testing
require "socket"
port = 8080
server = TCPServer.new port

system('bundle',  'install')
greeting = TomlRB.load_file('output.toml')['greeting']['output']

while session = server.accept
  request = session.gets
  puts request

  session.print "HTTP/1.1 200\r\n" # 1
  session.print "Content-Type: text/html\r\n" # 2
  session.print "\r\n" # 3
  session.print "#{greeting}! The time is #{Time.now}" #4

  session.close
end
