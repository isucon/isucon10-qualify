package Isuumo::Web;
use v5.32;
use warnings;
use utf8;
use Kossy;

use DBIx::Sunny;
use File::Spec;
use HTTP::Status qw/:constants/;
use Log::Minimal;
use JSON::MaybeXS;
use Cpanel::JSON::XS::Type;

local $Log::Minimal::LOG_LEVEL = "DEBUG";

our $MYSQL_CONNECTION_DATA = {
    host     => $ENV{MYSQL_HOST}     // '127.0.0.1',
    port     => $ENV{MYSQL_PORT}     // '3306',
    user     => $ENV{MYSQL_USER}     // 'isucon',
    dbname   => $ENV{MYSQL_DATABASE} // 'isuumo',
    password => $ENV{MYSQL_PASS}     // 'isucon',
};

sub dbh {
    my $self = shift;
    $self->{_dbh} ||= do {
        my ($host, $port, $user, $dbname, $password) = $MYSQL_CONNECTION_DATA->@{qw/host port user dbname password/};
        my $dsn = "dbi:mysql:database=$dbname;host=$host;port=$port";
        DBIx::Sunny->connect($dsn, $user, $password, {
            mysql_enable_utf8mb4 => 1,
            mysql_auto_reconnect => 1,
            Callbacks => {
                connected => sub {
                    my $dbh = shift;
                    # XXX $dbh->do('SET SESSION sql_mode="STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION"');
                    return;
                },
            },
        });
    };
}


use constant InitializeResponse => {
    language => JSON_TYPE_STRING
};

use constant ChairResponse => {
    id          => JSON_TYPE_INT,
    name        => JSON_TYPE_STRING,
    description => JSON_TYPE_STRING,
    thumbnail   => JSON_TYPE_STRING,
    price       => JSON_TYPE_INT,
    height      => JSON_TYPE_INT,
    width       => JSON_TYPE_INT,
    depth       => JSON_TYPE_INT,
    color       => JSON_TYPE_STRING,
    features    => JSON_TYPE_STRING,
    kind        => JSON_TYPE_STRING,
    popularity  => undef,
    stock       => undef,
};


post '/initialize' => sub {
    my ( $self, $c )  = @_;

    my $sql_dir = File::Spec->catdir($self->root_dir, "..", "mysql", "db");
    my @paths = (
        File::Spec->catfile($sql_dir, "0_Schema.sql"),
        File::Spec->catfile($sql_dir, "1_DummyEstateData.sql"),
        File::Spec->catfile($sql_dir, "2_DummyChairData.sql"),
    );

    for my $p (@paths) {
        my @cmd = ('mysql',
            '-h', $MYSQL_CONNECTION_DATA->{host},
            '-u', $MYSQL_CONNECTION_DATA->{user},
            "-p$MYSQL_CONNECTION_DATA->{password}",
            '-P', $MYSQL_CONNECTION_DATA->{port},
                  $MYSQL_CONNECTION_DATA->{dbname},
            '<',  $p);

        system(@cmd);
    }

    $self->res_json($c, {
        language => "perl",
    }, InitializeResponse);
};

get '/api/chair/{id}' => sub {
    my ( $self, $c )  = @_;

    my $chair_id = $c->args->{id};
    my $query = 'SELECT * FROM chair WHERE id = ?';
    my $chair = $self->dbh->select_row($query, $chair_id);

    if (!$chair) {
        infof("requested id's chair not found : %s", $chair_id);
        return $self->res_no_content($c, HTTP_NOT_FOUND)
    }

    if ($chair->{stock} <= 0) {
        infof("requested id's chair is sold out : %s", $chair_id);
        return $self->res_no_content($c, HTTP_NOT_FOUND)
    }

    return $self->res_json($c, {
        id          => $chair->{id},
        name        => $chair->{name},
        description => $chair->{description},
        thumbnail   => $chair->{thumbnail},
        price       => $chair->{price},
        height      => $chair->{height},
        width       => $chair->{width},
        depth       => $chair->{depth},
        color       => $chair->{color},
        features    => $chair->{features},
        kind        => $chair->{kind},
    }, ChairResponse)
};


# send empty body with status code
sub res_no_content {
    my ($self, $c, $status) = @_;
    $c->res->code($status);
    $c->res;
}

# render_json with json spec
# XXX: $json_specを指定できるようにKossy::Conection#render_jsonを調整
my $_JSON = JSON::MaybeXS->new()->allow_blessed(1)->convert_blessed(1)->ascii(1);
sub res_json {
    my ($self, $c, $obj, $json_spec) = @_;

    # defense from JSON hijacking
    # Copy from Amon2::Plugin::Web::JSON
    if ( exists $c->req->env->{'HTTP_X_REQUESTED_WITH'} &&
         ($c->req->env->{'HTTP_USER_AGENT'}||'') =~ /android/i &&
         exists $c->req->env->{'HTTP_COOKIE'} &&
         ($c->req->method||'GET') eq 'GET'
    ) {
        $c->halt(403,"Your request is maybe JSON hijacking.\nIf you are not a attacker, please add 'X-Requested-With' header to each request.");
    }

    my $body = $_JSON->encode($obj, $json_spec);
    $body = $c->escape_json($body);

    if ( ( $c->req->env->{'HTTP_USER_AGENT'} || '' ) =~ m/Safari/ ) {
        $body = "\xEF\xBB\xBF" . $body;
    }

    $c->res->status( 200 );
    $c->res->content_type('application/json; charset=UTF-8');
    $c->res->header( 'X-Content-Type-Options' => 'nosniff' ); # defense from XSS
    $c->res->body( $body );
    $c->res;
}

1;
