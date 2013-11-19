require 'rubygems'
require 'json'

default_run_options[:pty] = true
ssh_options[:forward_agent] = true


set :application, "GoAPI"
set :repository,  "git@github.com:curt-labs/GoAPI.git"
set :branch, "master"
set :deploy_via, :remote_cache

set :scm, :git
set :scm_passphrase, ""
set :user, "deployer"

role :web, "173.255.117.20", "173.255.112.170"
role :app, "173.255.117.20", "173.255.112.170"

set :deploy_to, "/home/#{user}/gocode/versionsing/#{application}"
set :app_path, "/home/#{user}/gocode/src/github.com/curt-labs/#{application}"
set :deploy_settings, "deploy_settings.json"

set :use_sudo, false
set :sudo_prompt, ""
set :normalize_asset_timestamps, false

set :default_environment, {
  'GOPATH' => '$HOME/gocode'
}

after "deploy", "deploy:goget"
after "deploy:goget", "db:configure"
after "db:configure", "deploy:compile"
after "deploy:compile", "deploy:stop"
after "deploy:stop", "deploy:restart"

namespace :db do
  desc "set database connction info"
  task :configure do
     if File.exists?(File.join(Dir.getwd, "config/#{deploy_settings}"))
      obj = JSON.parse(File.read(File.join(Dir.getwd, "config/#{deploy_settings}")))
      set(:database_host) { obj["database"]["host"] }
      set(:database_username) { obj["database"]["username"] }
      set(:database_password) { obj["database"]["password"]}
      set(:database_name) { obj["database"]["name"] }
    else
      set(:database_host) { Capistrano::CLI.ui.ask("Database Host: ") }
      set(:database_username) { Capistrano::CLI.ui.ask("Database Username: ") }
      set(:database_password) { Capistrano::CLI.password_prompt("Database Password: ")}
      set(:database_name) { Capistrano::CLI.ui.ask("Database Name: ") }
    end

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
      run "mkdir -p #{current_release}/current/helpers/database"
      put db_config, "#{current_release}/current/helpers/database/ConnectionString.go"
  end
end


namespace :deploy do
  task :goget do
  	run "/home/#{user}/bin/go get -u github.com/ziutek/mymysql/native"
  	run "/home/#{user}/bin/go get -u github.com/ziutek/mymysql/mysql"
  end
  task :compile do
    run "mkdir -p #{app_path}"
    run "cp -r #{current_release}/* #{app_path}"
  	run "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 /home/#{user}/bin/go build -o #{app_path}#{application} #{app_path}/index.go"
  end
  task :start do ; end
  task :stop do 
      kill_processes_matching "go-api"
      kill_processes_matching "#{application}"
  end
  task :restart do
    run "mkdir -p #{app_path}"
  	restart_cmd = "#{app_path}/#{application} -http=127.0.0.1:8080"
  	run "nohup sh -c '#{restart_cmd} &' > #{application}-nohup.out"
  end
end

def kill_processes_matching(name)
  run "ps -ef | grep #{name} | grep -v grep | awk '{print $2}' | sudo xargs kill -9 || echo 'no process with name #{name} found'"
end