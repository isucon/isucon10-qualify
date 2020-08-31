package Isuumo::Web;
use v5.32;
use warnings;
use utf8;
use Kossy;

use DBIx::Sunny;
use File::Spec;
use HTTP::Status qw/:constants/;
use Log::Minimal;

local $Log::Minimal::LOG_LEVEL = "DEBUG";

our $MYSQL_CONNECTION_DATA = {
    host     => $ENV{MYSQL_HOST}     // '127.0.0.1',
    port     => $ENV{MYSQL_PORT}     // '3306',
    user     => $ENV{MYSQL_USER}     // 'isucon',
    dbname   => $ENV{MYSQL_DATABASE} // 'isuumo',
    password => $ENV{MYSQL_PASS}     // 'isucon',
};

# send empty body with status code
sub res_no_content {
    my ($self, $c, $status) = @_;
    $c->res->code($status);
    $c->res;
}

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

    $c->render_json({
        language => "perl",
    });
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

    return $c->render_json($chair)
};


1;
