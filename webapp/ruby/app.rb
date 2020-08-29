require 'sinatra'
require 'mysql2'
require 'mysql2-cs-bind'
require 'csv'

class App < Sinatra::Base
  LIMIT = 20
  NAZOTTE_LIMIT = 50
  CHAIR_SEARCH_CONDITION = JSON.parse(File.read('../fixture/chair_condition.json'), symbolize_names: true)

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

  get '/api/chair/low_priced' do
    sql = "SELECT * FROM chair WHERE stock > 0 ORDER BY price ASC, id ASC LIMIT #{LIMIT}" # XXX:
    chairs = db.query(sql).to_a
    { chairs: chairs }.to_json
  end

  get '/api/chair/search' do
    search_queries = []
    query_params = []

    if params[:priceRangeId] && params[:priceRangeId].size > 0
      chair_price = CHAIR_SEARCH_CONDITION[:price][:ranges][params[:priceRangeId].to_i]
      unless chair_price
        logger.error "priceRangeID invalid: #{params[:priceRangeId]}"
        halt 400
      end

      if chair_price[:min] != -1
        search_queries << 'price >= ?'
        query_params << chair_price[:min]
      end

      if chair_price[:max] != -1
        search_queries << 'price < ?'
        query_params << chair_price[:max]
      end
    end

    if params[:heightRangeId] && params[:heightRangeId].size > 0
      chair_height = CHAIR_SEARCH_CONDITION[:height][:ranges][params[:heightRangeId].to_i]
      unless chair_height
        logger.error "heightRangeId invalid: #{params[:heightRangeId]}"
        halt 400
      end

      if chair_height[:min] != -1
        search_queries << 'height >= ?'
        query_params << chair_height[:min]
      end

      if chair_height[:max] != -1
        search_queries << 'height < ?'
        query_params << chair_height[:max]
      end
    end

    if params[:widthRangeId] && params[:widthRangeId].size > 0
      chair_width = CHAIR_SEARCH_CONDITION[:width][:ranges][params[:widthRangeId].to_i]
      unless chair_width
        logger.error "widthRangeId invalid: #{params[:widthRangeId]}"
        halt 400
      end

      if chair_width[:min] != -1
        search_queries << 'width >= ?'
        query_params << chair_width[:min]
      end

      if chair_width[:max] != -1
        search_queries << 'width < ?'
        query_params << chair_width[:max]
      end
    end

    if params[:depthRangeId] && params[:depthRangeId].size > 0
      chair_depth = CHAIR_SEARCH_CONDITION[:depth][:ranges][params[:depthRangeId].to_i]
      unless chair_depth
        logger.error "depthRangeId invalid: #{params[:depthRangeId]}"
        halt 400
      end

      if chair_depth[:min] != -1
        search_queries << 'depth >= ?'
        query_params << chair_depth[:min]
      end

      if chair_depth[:max] != -1
        search_queries << 'depth < ?'
        query_params << chair_depth[:max]
      end
    end

    if params[:kind] && params[:kind].size > 0
      search_queries << 'kind = ?'
      query_params << params[:kind]
    end

    if params[:color] && params[:color].size > 0
      search_queries << 'color = ?'
      query_params << params[:kind]
    end

    if params[:features] && params[:features].size > 0
      params[:features].split(',').each do |feature_condition|
        search_queries << "features LIKE CONCAT('%', ?, '%')"
        query_params.push(feature_condition)
      end
    end

    if search_queries.size == 0
      logger.error "Search condition not found"
      halt 400
    end

    search_queries.push('stock > 0')

    page =
      begin
        Integer(params[:page], 10)
      rescue ArgumentError => e
        logger.error "Invalid format page parameter: #{e.inspect}"
        halt 400
      end

    per_page =
      begin
        Integer(params[:perPage], 10)
      rescue ArgumentError => e
        logger.error "Invalid format perPage parameter: #{e.inspect}"
        halt 400
      end

    sqlprefix = 'SELECT * FROM chair WHERE '
    search_condition = search_queries.join(' AND ')
    limit_offset = " ORDER BY popularity DESC, id ASC LIMIT #{per_page} OFFSET #{per_page * page}" # XXX: mysql-cs-bind doesn't support escaping variables for limit and offset
    count_prefix = 'SELECT COUNT(*) as count FROM chair WHERE '

    count = db.xquery("#{count_prefix}#{search_condition}", query_params).first['count']
    chairs = db.xquery("#{sqlprefix}#{search_condition}#{limit_offset}", query_params).to_a

    { count: count, chairs: chairs }.to_json
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

  get '/api/chair/search/condition' do
    CHAIR_SEARCH_CONDITION.to_json
  end
end
