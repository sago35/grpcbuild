use strict;
use warnings;
use utf8;
use Term::Encoding qw(term_encoding);
my $encoding = term_encoding;

binmode STDOUT => "encoding($encoding)";
binmode STDERR => "encoding($encoding)";

use Path::Tiny;

use Getopt::Kingpin;
my $kingpin = Getopt::Kingpin->new;
my $in = $kingpin->arg("input", "input")->string;
my $right = $kingpin->arg("end", "end time (ms)")->default("0")->int;
$kingpin->parse();

my $p = path($in->value);

my $out = {};
{
    use bigint;
    my $start = 0;
    foreach my $line ($p->lines({chomp => 1})) {
        if ($line =~ /^#/) {
            my ($dum, $worker, $file, $time, $time2) = split /\s+/, $line;
            $start = $time if $start == 0;
            if (not defined $time2) {
                $out->{$worker}->{$file}->{start} = $time - $start;
            } else {
                $out->{$worker}->{$file}->{end} = $time - $start;
            }
        }
    }
    no bigint;
}


my $str = <<'...';
<html>
<head>
<style>
rect{ stroke:black; stroke-width:1; }
</style>
</head>
<body>
<script type="text/javascript" src="https://www.gstatic.com/charts/loader.js"></script>

<script type="text/javascript">
  google.charts.load("current", {packages:["timeline"]});
  google.charts.setOnLoadCallback(drawChart);
  function drawChart() {

    var container = document.getElementById('example3.1');
    var chart = new google.visualization.Timeline(container);
    var dataTable = new google.visualization.DataTable();
    dataTable.addColumn({ type: 'string', id: 'Position' });
    dataTable.addColumn({ type: 'string', id: 'Name' });
    dataTable.addColumn({ type: 'date', id: 'Start' });
    dataTable.addColumn({ type: 'date', id: 'End' });
    dataTable.addRows([
...

my $end = {
    end => 0,
};

foreach my $worker (sort mysort keys %$out) {
    foreach my $file (keys %{$out->{$worker}}) {
        my $t = $out->{$worker}->{$file};
        my $f = $file;
        $f =~ s/^testdata\///;
        $f =~ s/\\/\//g;
        $str .= sprintf "    ['%s', '%s', new Date(%s), new Date(%s)],\n", $worker, $f, $t->{start}, $t->{end};
        if ($end->{end} < $t->{end}) {
            $end = {
                worker => $worker,
                file   => "",
                start  => $t->{end},
                end    => $t->{end},
            };
        }
    }
}

if ($end->{end} < $right->value) {
    $str .= sprintf "    ['%s', '%s', new Date(%s), new Date(%s)],\n", $end->{worker}, $end->{file}, $right->value, $right->value;
}

$str =~ s/,$//;

$str .= <<'...';
    ]);

    var options = {
      timeline: {
        avoidOverlappingGridLines: false,
        colors: ['#cbb69d', '#603913', '#c69c6e'],
        showBarLabels: false,
        colorByRowLabel: true,
        groupByRowLabel: true
      }
    };

    chart.draw(dataTable, options);
  }
</script>

<div id="example3.1" style="height: 2200px;"></div>
<body>
</html>
...

print $str;

sub mysort {
    my @aa = split /[\(\)]/, $a;
    my @bb = split /[\(\)]/, $b;

    if ($aa[0] ne $bb[0]) {
        return $bb[0] cmp $aa[0];
    } elsif ($aa[1] ne $bb[1]) {
        return $aa[1] <=> $bb[1];
    }

}
