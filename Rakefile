task :test do
    sh "go test -timeout 10s ./..."
end

task :build do
    sh "GOOS=linux GOARCH=amd64 go build -o gojikoServer.linux.amd64 goa/*.go"
end
