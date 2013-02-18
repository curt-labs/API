set :application, "go-api"
set :repository,  "git@github.com:curt-labs/GoAPI.git"

ssh_options[:forward_agent] = true

set :user, 'ninnemana'
set :deploy_to, "/home/#{user}/app"
set :gopath, deploy_to
set :pid_file, deploy_to+'/pids/PIDFILE'
set :symlinks, { "pids" => "pids"}

role :app, "curt-api.cloudapp.net"
role :web, "curt-api.cloudapp.net"

task :production do
	server "curt-api.cloudapp.net", :app
end

after 'deploy:update_code', 'go:build'

namespace :go do
	task :build do
		run "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o go-api.linux index.go"
		run "scp go-api.linux deploy-user:eC0mm3rc3@curt-api.cloudapp.net:/home/ninnemana"
	end
end
