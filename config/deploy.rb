default_run_options[:pty] = true
ssh_options[:forward_agent] = true


set :application, "GoAPI"
set :repository,  "git@github.com:curt-labs/GoAPI.git"

set :scm, :git
set :scm_passphrase, ""
set :user, "ninnemana"

role :web, "curt-api-server1.cloudapp.net", "curt-api-server2.cloudapp.net"
role :app, "curt-api-server1.cloudapp.net", "curt-api-server2.cloudapp.net"

set :deploy_to, "/home/ninnemana/#{application}"
set :deploy_via, :remote_cache

set :use_sudo, false
set :sudo_prompt, ""
set :normalize_asset_timestamps, false

after "deploy", "deploy:goget"
after "deploy:goget", "deploy:compile"
after "deploy:compile", "deploy:stop"
after "deploy:stop", "deploy:restart"

namespace :deploy do
  task :goget do
  	run "/usr/local/go/bin/go get github.com/ziutek/mymysql/native"
  	run "/usr/local/go/bin/go get github.com/ziutek/mymysql/mysql"
  end
  task :compile do
  	run "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 /usr/local/go/bin/go build -o #{deploy_to}/current/go-api #{deploy_to}/current/index.go"
  end
  task :start do ; end
  task :stop do 
      kill_processes_matching "go-api"
  end
  task :restart do
  	restart_cmd = "#{current_release}/go-api"
  	run "nohup sh -c '#{restart_cmd} &' > nohup.out"
  end
end

def kill_processes_matching(name)
  run "ps -ef | grep #{name} | grep -v grep | awk '{print $2}' | sudo xargs kill -2 || echo 'no process with name #{name} found'"
end