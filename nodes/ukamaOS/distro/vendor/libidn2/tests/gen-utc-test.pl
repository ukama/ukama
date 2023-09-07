#!/usr/bin/perl

# Copyright (C) 2011-2021 Simon Josefsson

# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU General Public License for more details.
#
# You should have received a copy of the GNU General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

# I consider the output of this program to be unrestricted.  Use it as
# you will.

use strict;

my ($last);
my ($lineno) = 0;
my ($ctr) = 0;

while (<>) {
    $lineno++;
    next unless /^[BN]/;

    next unless m,^.*;\t(.*);\t(.*);\t(.*);\s(NV8).*,;

    my $line = $_;

    my ($source) = $1;
    my ($ustr) = $2;
    my ($astr) = $3;
    my ($nv8) = $4;

    $ustr = $source if ($ustr eq "");
    $astr = $ustr if ($astr eq "");

    while ($ustr =~ /(.*)\\u([0-9A-f][0-9A-f][0-9A-f][0-9A-f])(.*)/) {
	my $num = hex($2);
	#printf "/* hex $2 num $num */";

	my $str = unpack ("H*", pack("C0U*",$num));
	my $escstr = "";
	while ($str) {
	    $escstr .= "\\x" . substr ($str,0,2);
	    $str = substr ($str,2);
	}
	#printf "/* utf8 $escstr */\n";

	$ustr = $1.'" "'.$escstr.'" "'.$3;
    }

    next if ($ustr eq $last);

    print "/* lineno $lineno ctr $ctr source $source uni $ustr ace $astr nv8 $nv8 line $line */\n";

    if ($astr =~ /\\u/) {
	print "/* IdnaTest.txt bug? */\n";
    } elsif ($astr =~ /。/) {
	print "/* IdnaTest.txt bug2? */\n";
    } elsif ($ustr =~ /a..c/ || $ustr =~ /ä..c/) {
	print "/* libidn2 bug? */\n";
    } elsif ($nv8 eq "NV8") {
	print "{ \"$ustr\", \"$astr\", -1 },\n";
	$ctr++;
    } elsif (substr($astr, 0, 1) eq "[" && substr($ustr, 0, 1) ne "[") {
	print "{ \"$ustr\", \"$astr\", -1 },\n";
	$ctr++;
    } elsif (substr($astr, 0, 1) eq "[") {
	print "/* punt1 $line */\n";
    } else {
	print "{ \"$ustr\", \"$astr\", IDN2_OK },\n";
	$ctr++;
    }

    $last = $ustr;
}
