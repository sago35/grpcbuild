use strict;
use warnings;
use utf8;
use Term::Encoding qw(term_encoding);
my $encoding = term_encoding;

binmode STDOUT => "encoding($encoding)";
binmode STDERR => "encoding($encoding)";

use Path::Tiny;
my $t = path("./testdata");

srand(time);

foreach my $x ('a' .. 'f') {
    foreach my $y ('a' .. 'f') {
        my $wait_ms = int(rand(2000) + 200);
        my $lines = int(rand(20)) - 5;
        $lines = $lines < 0 ? 0 : $lines;

        printf "%s %s %d %d\n", $x, $y, $wait_ms, $lines;

        my @data = (sprintf "%d\n", $wait_ms);
        for (my $i = 0; $i < $lines; $i++) {
            push @data, sprintf "warning %s%s%02d\n", $x, $y, $i + 1;
        }

        my $f = $t->child(sprintf "%s%s.c", $x, $y);
        $f->spew(@data);
    }
}
