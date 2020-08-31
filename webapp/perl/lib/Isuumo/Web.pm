package Isuumo::Web;

use strict;
use warnings;
use utf8;
use Kossy;

use File::Spec;

our $MYSQL_CONNECTION_DATA = {
    host     => $ENV{MYSQL_HOST}   // '127.0.0.1',
    port     => $ENV{MYSQL_PORT}   // '3306',
    user     => $ENV{MYSQL_USER}   // 'isucon',
    dbname   => $ENV{MYSQL_DBNAME} // 'isuumo',
    password => $ENV{MYSQL_PASS}   // 'isucon',
};

get '/initialize' => sub {
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

1;
