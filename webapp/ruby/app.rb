require 'sinatra'
require 'mysql2'
require 'mysql2-cs-bind'
require 'csv'

class App < Sinatra::Base
  LIMIT = 20
  NAZOTTE_LIMIT = 50

  configure :development do
    require 'sinatra/reloader'
    register Sinatra::Reloader
  end

  configure do
    enable :logging
  end

  set :add_charset, ['application/json']

  helpers do
    def db
      Thread.current[:db] ||= Mysql2::Client.new(
        host: ENV.fetch('MYSQL_HOST', '127.0.0.1'),
        port: ENV.fetch('MYSQL_PORT', '3306'),
        username: ENV.fetch('MYSQL_USER', 'isucon'),
        password: ENV.fetch('MYSQL_PASS', 'isucon'),
        database: ENV.fetch('MYSQL_DBNAME', 'isuumo'),
        reconnect: true,
      )
    end
  end

  post '/initialize' do
    unless system('../mysql/db/init.sh')
      logger.error('Initialize script error')
      halt 500
    end

    { language: 'ruby' }.to_json
  end

  get '/api/chair/:id' do
    id =
      begin
        Integer(params[:id], 10)
      rescue ArgumentError => e
        logger.error "Request parameter \"id\" parse error: #{e.inspect}"
        halt 400
      end

    chair = db.xquery('SELECT * FROM chair WHERE id = ?', id).first
    unless chair
      logger.info "Requested id's chair not found: #{id}"
      halt 404
    end

    if chair['stock'] <= 0
      logger.info("Requested id's chair is sold out: #{id}")
      halt 404
    end

    chair.to_json
  end

  post '/api/chair' do
    if !params[:chairs] || !params[:chairs].respond_to?(:key) || !params[:chairs].key?(:tempfile)
      logger.error 'Failed to get form file'
      halt 400
    end

    db.query('BEGIN')
    begin
      CSV.parse(params[:chairs][:tempfile].read) do |row|
        sql = 'INSERT INTO chair(id, name, description, thumbnail, price, height, width, depth, color, features, kind, popularity, stock) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)'
        db.xquery(sql, *row)
      end
    rescue => e
      db.query('ROLLBACK')
      logger.error("Failed to commit tx: #{e.inspect}")
      halt 500
    end

    db.query('COMMIT')

    status 201
  end
end
