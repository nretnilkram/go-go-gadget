package strtwist

import (
	"testing"
)

func TestSymbolSubstitution(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"abcdefghijklmnopqrstuvwxyz", "@bcd3fgh!jklmn0pqr$tuvwxyz"},
		{"ABCDEFGHIJKLMNOPQRSTUVWXYZ", "@BCD3FGH!JKLMN0PQR$TUVWXYZ"},
		{"ABCDEF GHIJKLMNOP QRSTUVW XYZ", "@BCD3F GH!JKLMN0P QR$TUVW XYZ"},
		{"partner pickle steam violet-blue act fit birds marten blue-violet few shamrock lip dust coyote pine-green hope nail yellow-orange net angle cerulean important moon sorry lovebird wild-watermelon reindeer tumbleweed retain old wild-blue-yonder care burnt-sienna wild-blue-yonder dandelion payment shocking-pink wombat view horror such rich chimpanzee crate cactus gold catch shot tiny cow ferret tap suck airline finch main radical-red orange-yellow antique-brass daughter nearby choice voice district vivid-violet tropical-rain-forest easy is maximum feature koodoo former flag royal-purple forest-green promise skunk mountain turquoise-blue turkey interesting exchange fixed blush mahogany united weather unhappy lavender guitar former settle history popcorn daughter measure distinct smile electric-lime mirror", "p@rtn3r p!ckl3 $t3@m v!0l3t-blu3 @ct f!t b!rd$ m@rt3n blu3-v!0l3t f3w $h@mr0ck l!p du$t c0y0t3 p!n3-gr33n h0p3 n@!l y3ll0w-0r@ng3 n3t @ngl3 c3rul3@n !mp0rt@nt m00n $0rry l0v3b!rd w!ld-w@t3rm3l0n r3!nd33r tumbl3w33d r3t@!n 0ld w!ld-blu3-y0nd3r c@r3 burnt-$!3nn@ w!ld-blu3-y0nd3r d@nd3l!0n p@ym3nt $h0ck!ng-p!nk w0mb@t v!3w h0rr0r $uch r!ch ch!mp@nz33 cr@t3 c@ctu$ g0ld c@tch $h0t t!ny c0w f3rr3t t@p $uck @!rl!n3 f!nch m@!n r@d!c@l-r3d 0r@ng3-y3ll0w @nt!qu3-br@$$ d@ught3r n3@rby ch0!c3 v0!c3 d!$tr!ct v!v!d-v!0l3t tr0p!c@l-r@!n-f0r3$t 3@$y !$ m@x!mum f3@tur3 k00d00 f0rm3r fl@g r0y@l-purpl3 f0r3$t-gr33n pr0m!$3 $kunk m0unt@!n turqu0!$3-blu3 turk3y !nt3r3$t!ng 3xch@ng3 f!x3d blu$h m@h0g@ny un!t3d w3@th3r unh@ppy l@v3nd3r gu!t@r f0rm3r $3ttl3 h!$t0ry p0pc0rn d@ught3r m3@$ur3 d!$t!nct $m!l3 3l3ctr!c-l!m3 m!rr0r"},
	}
	for _, c := range cases {
		got := SymbolSubstitution(c.in)
		if got != c.want {
			t.Errorf("K8s(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}
