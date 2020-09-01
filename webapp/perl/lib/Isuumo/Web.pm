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
use Text::CSV_XS;

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

use constant Chair => {
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

use constant ChairSearchResponse => {
    count  => JSON_TYPE_INT,
    chairs => json_type_arrayof(Chair),
};

use constant ChairListResponse => {
    chairs => json_type_arrayof(Chair),
};

# Estate 物件
use constant Estate => {
    id          => JSON_TYPE_INT,
    thumbnail   => JSON_TYPE_STRING,
    name        => JSON_TYPE_STRING,
    description => JSON_TYPE_STRING,
    latitude    => JSON_TYPE_FLOAT,
    longitude   => JSON_TYPE_FLOAT,
    address     => JSON_TYPE_STRING,
    rent        => JSON_TYPE_INT,
    doorHeight  => JSON_TYPE_INT,
    doorWidth   => JSON_TYPE_INT,
    features    => JSON_TYPE_STRING,
    popularity  => undef,
};

use constant EstateSearchResponse => {
    count   => JSON_TYPE_INT,
    estates => json_type_arrayof(Estate),
};

use constant EstateListResponse => {
    estates => json_type_arrayof(Estate),
};

use constant Coordinate => {
    latitude    => JSON_TYPE_FLOAT,
    longitude   => JSON_TYPE_FLOAT,
};

use constant Coordinates => {
    coordinates => json_type_arrayof(Coordinate),
};

use constant Range => {
    id  => JSON_TYPE_INT,
    min => JSON_TYPE_INT,
    max => JSON_TYPE_INT,
};

use constant RangeCondition => {
    prefix => JSON_TYPE_STRING,
    suffix => JSON_TYPE_STRING,
    ranges => json_type_arrayof(Range),
};

use constant ListCondition => {
    list => json_type_arrayof(JSON_TYPE_STRING),
};

use constant EstateSearchCondition => {
    doorWidth  => RangeCondition,
    doorHeight => RangeCondition,
    rent       => RangeCondition,
    feature    => ListCondition,
};

use constant ChairSearchCondition => {
    width   => RangeCondition,
    height  => RangeCondition,
    depth   => RangeCondition,
    price   => RangeCondition,
    color   => ListCondition,
    feature => ListCondition,
    kind    => ListCondition,
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
    }, Chair)
};

post '/api/chair' => sub {
    my ( $self, $c )  = @_;

    my $file = $c->req->uploads->{'chairs'};
    if (!$file) {
        critf("failed to get form file");
        return $self->res_no_content($c, HTTP_BAD_REQUEST);
    }

    my $fh;
    if (!open $fh, "<:encoding(utf8)", $file->path) {
        critf("failed to open form file %s. %s", $file->path, $!);
        return $self->res_no_content($c, HTTP_INTERNAL_SERVER_ERROR);
    }

    my $csv = Text::CSV_XS->new({binary => 1});
    my $dbh = $self->dbh;
    my $txn = $dbh->txn_scope;

    eval {
        while (my $row = $csv->getline($fh)) {
            my ($id, $name, $description, $thumbnail, $price, $height, $width, $depth, $color, $features, $kind, $popularity, $stock) = $row->@*;
            $dbh->query(
                "INSERT INTO chair(id, name, description, thumbnail, price, height, width, depth, color, features, kind, popularity, stock) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?)",
                $id, $name, $description, $thumbnail, $price, $height, $width, $depth, $color, $features, $kind, $popularity, $stock
            );
        }
        $txn->commit;
    };

    $fh->close;

    if ($@) {
        $txn->rollback;
        critf("failed to commit txn: %s", $@);
        return $self->res_no_content($c, HTTP_INTERNAL_SERVER_ERROR);
    }
    return $self->res_no_content($c, HTTP_CREATED);
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
