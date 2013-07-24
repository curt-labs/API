default_run_options[:pty] = true
ssh_options[:forward_agent] = true


set :application, "GoAPI"
set :repository,  "git@github.com:curt-labs/GoAPI.git"

set :scm, :git
set :scm_passphrase, ""
set :user, "deployer"

# role :web, "curt-api-server1.cloudapp.net", "curt-api-server2.cloudapp.net"
# role :app, "curt-api-server1.cloudapp.net", "curt-api-server2.cloudapp.net"
role :web, "173.255.117.20", "173.255.112.170"
role :app, "173.255.117.20", "173.255.112.170"

set :deploy_to, "/home/#{user}/#{application}"
set :deploy_via, :remote_cache
set :gopath, deploy_to

set :use_sudo, false
set :sudo_prompt, ""
set :normalize_asset_timestamps, false


after "deploy", "deploy:goget"
after "deploy:goget", "db:configure"
after "db:configure", "deploy:compile"
after "deploy:compile", "deploy:stop"
after "deploy:stop", "deploy:restart"


namespace :db do
  desc "set database connction info"
  task :configure do
    set(:database_host) { Capistrano::CLI.ui.ask("Database Host: ") }
    set(:database_username) { Capistrano::CLI.ui.ask("Database Username: ") }
    set(:database_password) { Capistrano::CLI.password_prompt("Database Password: ")}
    set(:database_name) { Capistrano::CLI.ui.ask("Database Name: ") }

    db_config = <<-EOF
      package database

      const (
        db_proto = "tcp"
        db_addr = "#{database_host}"
        db_user = "#{database_username}"
        db_pass = "#{database_password}"
        db_name = "#{database_name}"
      )
      EOF
      run "mkdir -p #{deploy_to}/current/helpers/database"
      put db_config, "#{deploy_to}/current/helpers/database/ConnectionString.go"
  end
end


namespace :deploy do
  task :goget do
    with_env('GOPATH', gopath) do
    	#run "/home/#{user}/bin/go get -u github.com/ziutek/mymysql/native"
    	# run "export GOPATH=$HOME/gocode | sudo /home/#{user}/bin/go get -u github.com/ziutek/mymysql/mysql"
    end
  end
  task :compile do
  	run "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 /home/#{user}/bin/go build -o #{deploy_to}/current/go-api #{deploy_to}/current/index.go"
  end
  task :start do ; end
  task :stop do 
      kill_processes_matching "go-api"
  end
  task :restart do
  	restart_cmd = "#{current_release}/go-api -http=127.0.0.1:8080"
  	run "nohup sh -c '#{restart_cmd} &' > #{application}-nohup.out"
  end
end

def kill_processes_matching(name)
  run "ps -ef | grep #{name} | grep -v grep | awk '{print $2}' | sudo xargs kill -9 || echo 'no process with name #{name} found'"
end