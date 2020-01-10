package ad

import (
	"encoding/json"
	"fmt"
	sdk "github.com/irisnet/irishub/types"
	"github.com/irisnet/sdk-go/client/basic"
	"github.com/irisnet/sdk-go/client/lcd"
	"github.com/irisnet/sdk-go/client/rpc"
	"github.com/irisnet/sdk-go/client/tx"
	ctypes "github.com/irisnet/sdk-go/client/types"
	"github.com/irisnet/sdk-go/keys"
	"github.com/irisnet/sdk-go/types"
	"github.com/irisnet/sdk-go/util"
	"math"
	"testing"
	"time"
)

type (
	MnemonicKMParams struct {
		Menmonic string
		Password string
		FullPath string
	}

	KeyStoreKMParams struct {
		FilePath string
		Password string
	}
)

var (
	liteClient lcd.LiteClient
	rpcClient  rpc.RPCClient
	txClient   tx.TxClient
	km         keys.KeyManager
)

func initClient() {
	network := types.Mainnet
	lcdUrl := "http://v2.irisnet-lcd.rainbow.one"
	rpcUrl := "tcp://seed-1.mainnet.irisnet.org:26657"
	initKMType := "seed"

	//network := types.Testnet
	//lcdUrl := "http://irisnet-lcd.dev.bianjie.ai"
	//rpcUrl := "tcp://192.168.150.32:26657"
	//initKMType := "seed"

	switch initKMType {
	case "seed":
		p := MnemonicKMParams{
			Menmonic: "proof own domain feature vehicle excite exotic way monkey stuff animal gorilla security roast street artwork room blue smoke fancy address",
			Password: "",
			FullPath: "44'/118'/0'/0/0",
		}
		if v, err := keys.NewMnemonicKeyManager(p.Menmonic, p.Password, p.FullPath); err != nil {
			fmt.Println("init km fail")
			panic(err)
		} else {
			km = v
		}
		break
	case "ks":
		p := KeyStoreKMParams{
			FilePath: "",
			Password: "",
		}
		if v, err := keys.NewKeyStoreKeyManager(p.FilePath, p.Password); err != nil {
			fmt.Println("init km fail")
			panic(err)
		} else {
			km = v
		}
		break
	default:
		panic("should init km first")
	}

	switch network {
	case types.Mainnet:
		sdk.SetNetworkType("mainnet")
	case types.Testnet:
		sdk.SetNetworkType("testnet")
	}

	fmt.Printf("address is: %s\n", km.GetAddr().String())

	basicClient := basic.NewClient(lcdUrl)
	liteClient = lcd.NewClient(basicClient)
	rpcClient = rpc.NewClient(rpcUrl)
	var chainId string
	if status, err := rpcClient.GetStatus(); err != nil {
		fmt.Println("init rpc client err")
		panic(err)
	} else {
		chainId = status.NodeInfo.Network
	}

	if c, err := tx.NewClient(chainId, network, km, liteClient, rpcClient); err != nil {
		fmt.Println("init tx client err")
		panic(err)
	} else {
		txClient = c
	}
}

func TestInitClient(m *testing.T) {
	initClient()
}

func TestGetAccountInfo(t *testing.T) {
	initClient()
	address := "iaa179f42ymtd6ltslm3h6tlchs8s3ldwaj6lvjhfv"
	if v, err := liteClient.QueryAccount(address); err != nil {
		t.Fatal(err)
	} else {
		t.Log(util.ToJsonIgnoreErr(v))
	}
}

func TestSendToken(t *testing.T) {
	initClient()
	memo := "Rainbow airdrop for experience IRISnet proposal voting"
	addrAmountMap := map[string]float64{}

	//jsonDataStr := `{"iaa1qtm7hw464wkr473k8qu30hgvh22gatdkc3w7p6":0.02,"iaa1j7xaagn3ummuxwt9s9lwdykdk48clyvy85gcmr":0.01}`

	jsonDataStr := `{"iaa1026qrwaaejjavscnvepjm5t670n2s0ucffhed0":6,"iaa102q4gpvh4kmaus402lgcn4u2ylfx80wr6e3ryt":6,"iaa10cp5qa7hn6fj2xtaz7530n2jzfa7qmnsg2qdf9":6,"iaa10gfg38wyn9r83uzrh4z3fytfz3qsxdtpdcq2mh":6,"iaa10l0khu6ddhwlej94eptedznp7mj40g2enpe0s8":6,"iaa10mxldppuzcxjt23c07mtzx2xz3whc56ycrej2l":6,"iaa10padu9n2q7xqhdkde3umjtwpdrj9f3g5r4pjl6":6,"iaa10qthtarg3mcgpezvr74gqaxv8hdms53mcw9a4x":6,"iaa10sp709q23q6kfgwcpc8syvf9culywr2y3wg2c4":6,"iaa10tp6zc4lvcku5cfzltxl7y7vtwmuj9htg8yy5r":6,"iaa10v9kf6u70rtqzew9l0a58j0xvq76r2zucmj78q":6,"iaa10x0mva68vvqesgu4hd94ayfml0a382y53fl9hd":6,"iaa124y3raq3flwt9n8e3qwjate9n2u2mrmep64yk0":6,"iaa12c8ckm9su3jamy9czqe7ysyfl49qk5p8af5dah":6,"iaa12fd49g2c8vf5gv8ljz6aut397knnnjmjgk8axe":6,"iaa12l9h8jmrymtdchejns3huegnern08rrel3kfr9":6,"iaa12psd3exwxs2zgy9qrrmc37zn8pv2au8yjhs5lz":6,"iaa12s0ddw23glasqh8e3pjpx40vcy3q9fhuqfq9xl":6,"iaa12u3g03ty5k6tplwek0zqknvcg6qtmslku26u05":6,"iaa12ufzlf78ptu8s9sn4lrwvertesefyhvvnacc65":6,"iaa12uu9j4v2p3clec429c9dxnyjk7z8v62cy5wtg5":6,"iaa1307r6skc0d6qkl6r50kgxxeases0dh2zh0tllw":6,"iaa13248jes6rrswdvz8lljsda7dfelynsyjqrfcz0":6,"iaa132yd5q8vdmrnjz3ksdnmqfpx9q0zn8weq73edw":6,"iaa1330008yykgwv07wk8te3a3wvjea6l8ypqddfxl":6,"iaa1338yeawvta4a4zgtysez068u56met2k3pckp2u":6,"iaa138kla03j07xj2rln5p3v5xvjrx8vvyk62y0rmw":6,"iaa13dur05xudrpnmgukuca6zwluyn70kx39qcqg52":6,"iaa13e4rxrqgcgjwrnyeq99cz6run59tc794yxwrq4":6,"iaa13kk0xtedkunyl2a3u46p6um4ae0dz2rjatcye9":6,"iaa13quvlsugg6jcsur55nh9njt82luvsnpg94ek5j":6,"iaa13tswxx3vzfg7apfh2l0jl54z87jlz6pq7ppw50":6,"iaa140qzvm0g6v4et8wys53rcyy69huddh6savd5jk":6,"iaa143suv89arrf6786gyu23446cxq7eysege7522a":6,"iaa1484ar040gd7ppq0rf667u7ptx8n0u6z4aswjkd":6,"iaa14ha5urqqlk9phh4wmzjh5v9zldy663r8ys87ts":6,"iaa14ltzaah25djm49g4hqfcg3zujad8dzy66yfd6s":6,"iaa14nhqjude62xjsmv0wa7cagpwrvp293ddsuaeec":6,"iaa14u945vxnjhh8umnq8n3kvtay3em24dzmj05nvc":6,"iaa14xtmx49cl7l92eqr5tz54taj8slzte25r2pqg9":6,"iaa15gqfv3nmy4hqhj3hk0hlgdjlmejstewkygzwsz":6,"iaa15h7awla6rem8ynj85p5wx8qujpcztmvppday3q":6,"iaa15u9lkx734zcy6g5v8kkzr439ugar7ja2s0r86d":6,"iaa1654xshv4zyxarhtafw3483nfw04j8vs9g7pc9f":6,"iaa167ug37cd6dmm0s0jdgjwpe0qwjhenh5u9za60k":6,"iaa16hvxxftn9lkt5fwdsamkn7fksy5d6zg8f4w5mr":6,"iaa16qf9pn6u8kp2krzmw9g4qvw7l5ykamps0dezve":6,"iaa16xf95cpuk979p0d7tjterm09jgglscrx0xrjl3":6,"iaa178jznkvsvqemrwmyz7cf4tvdv9cv6xnamp3fea":6,"iaa17dftdlzrzt9e09hmzr5cfs2n9zth52gyu080n0":6,"iaa17efdjm76myzfxw479eqwxsh9csc2565flgls2k":6,"iaa17l0pwy4qvplz5s0zjg2tgp995999kna730fmts":6,"iaa17qhy6cuqyz4g4xefwaftzwgdse4va3z8y65fqa":6,"iaa17v9q69ax8jftlgwuxn5mqakjz2jf230qkseqtf":6,"iaa17vyu0ewjl0kkzuczv795wck94julpw3xhl2l4v":6,"iaa180qcntxnh6lj0a6jmxfcpjf32k97u6l9xs8nnw":6,"iaa185vfkxw6cl6samd4thnw2g0kk7ueh0q4rhn27r":6,"iaa188q82c2v6ehc88r8ek0fn9mf3u5z92lcst9ulk":6,"iaa188zwzfvjqhc7fakyj7yw8d4fq24agmv0yadvzm":6,"iaa18dvz99edrp3x6t4maxmyjnu02eg7r28zha5ur6":6,"iaa18fu2jj3ns95kjhjtz85gwzskqsr4hz47z904t2":6,"iaa18n37ez6qp6sp6asacpw3cqw3sgv9033m2r7anh":6,"iaa18pxx4297wenyu7awgfy3nwghzgsvk5pc58hl35":6,"iaa18skv95m9qqjklhtzzekdck4wht8ejt5k8qply4":6,"iaa18v2z7er8mls0ceceu82jwefnjdvnw4xkmwlv75":6,"iaa18xgeau6ehq2qg2khh72egg9hq6xgtvvt7wmre0":6,"iaa18yr3jfdrkz58x6xpus9rxq35x8mq0wjp5w6trh":6,"iaa193t9ejjr63n0djw5wdkkvl0vqkuafew053ca90":6,"iaa199xpvkx7sanrgqtufjetqjd9cknpjq9af4lrye":6,"iaa19atq5qj3snp000lz0mn40d7d2tkp350pkuu39m":6,"iaa19ch2jyl5v4jgazeprfv6v7jj6d70d6qmj6ex55":6,"iaa19cmhg32kkc668sjsnakuvga5kq4d02prf5f6yu":6,"iaa19d3kddcyz5ekwj9xnzg85met8damr80fgnelyh":6,"iaa19qdhqcu56x09jjs6e79pp35temq0hm0h4hfn6r":6,"iaa19rahjgjq9n962l3ezkpa3vm6g898hatwax3kap":6,"iaa19tk96rvgqxe604m6zypyl7fuupkp4pckyayaeh":6,"iaa19tuuj9vj9n9jlrjyn4tp65ay7w5vptsnfwwdlk":6,"iaa19vwdxqkwrtnd4h7z5e9zz3qwa06jd0wn5n3qf9":6,"iaa19yjtks2wwrn0fqusku68es9l336e5z3gl62lne":6,"iaa1a3dz7prw9q8jr54czscay26wnuw2knsfp6u0e9":6,"iaa1a8n8qk8dl8rdx6pem8wa6st8fhefnyy0qy424u":6,"iaa1a9g9hcwven940qwv8z2v5cpl3sneu9pzvze82y":6,"iaa1aepg660fyg7lzwzkk2388zj80k69e4wdwn43qv":6,"iaa1c5ckcxlqdqrzwl0932kcmx82z90fu3q8gndzg3":6,"iaa1c6eh5zca9cv0495a8juas47hsyp5whn89zy529":6,"iaa1c88fm7faeqn59fvjwafuavdrmcwakv73fcyupf":6,"iaa1c9ptf0gewl9xjvrjw0ct9wtjk8q9mhcjpp4x2z":6,"iaa1cdprulu3m5hk99el8c9wc4d208srrm627wcdr4":6,"iaa1cmd640lg0wu96s05rthzjt38agntpg5wz7x958":6,"iaa1cmqgy38y62xplcsa4cqcucatretgdt66j0v9h4":6,"iaa1cmxh0zjqy9uq8c0vvd27z4x2cg7aajp7scrkh7":6,"iaa1cnchyft0juwjkrj4cch9mxjdrq9zp8z9mtulmc":6,"iaa1cp04ymk4f08283fe4dqaaxfljvuurpnlqgr3ka":6,"iaa1cp9euwryr6mqlkufvxlmafmdatd0kwa4vgct5z":6,"iaa1cwypr8kslyr6naevr9s2djw85qk2dvsgrxkzdw":6,"iaa1cyd8pmrdwllqrnavlklvr48qpm38zndnrn5gsa":6,"iaa1d0d3qf7kc8jct4zvsfuvpallanc5vkpa2ms0m5":6,"iaa1d8ecgh2j9nrrzashs6geeftgq36qd7hdsnrql5":6,"iaa1d8k6qsgjt4ywytwj0qv52pg8h7y8j6uujankvs":6,"iaa1ddckxhypgpdz0hnckd4evphqe0sqdscjnv0jxn":6,"iaa1dkg9j2mmlmwg7uslaxc0hvn8e30grkedcte43s":6,"iaa1e33w7hmanmtsldwgzqlkfj6kwmxljs33kjkawa":6,"iaa1e5zrx4yus4nkwy2v4utrz3uvrngan0pjmnhzu4":6,"iaa1e9l6fhwl7jt0sddl63smprt34xsf5lszdc2g6n":6,"iaa1eletxstz8aup6p8uvq0ku8lkg5999hzve0v3vf":6,"iaa1em29a6gvud976nccc4gj8u7txr04z24t3rlg5e":6,"iaa1epq0wzld73c2vt3xngrak973mgth8fuazcw580":6,"iaa1exrtdfy8czgfm8k2cuajdcyyfzecaecm99saqk":6,"iaa1f93zjgn6dc6gl2yqjsh63cu6ykup63399gmy7c":6,"iaa1fchk7ahkmavsr7jscrua93w65rr5t4qyhuzduz":6,"iaa1ff0nyptye6esfx9tzk539knkkxxtcape86vwdt":6,"iaa1fga5rw2mcr4wlduezlqk4ypmx53fqy32q2ngm0":6,"iaa1fke62y5ds9hkmnqkc0409qkz3dl4vj4jks2vrl":6,"iaa1fsgap9jdthlxznfy2zdv966xtla59at6ug59t8":6,"iaa1ft747dwuc9hn6xltmqmpmt4tdqs2ck37qgpjjs":6,"iaa1g22efp6lln6jqw5yhhjy2q5kdjhuphx46rdygy":6,"iaa1g3hsvqmp72lkuru33m4d9k292jlu7329e7fefr":6,"iaa1gg6auuyv4h4pxfwsv2rj8p4myt0xj5a6l9x2ml":6,"iaa1glyjp4muty76g0rms9c5dx6hwhfmvtgpld2dvw":6,"iaa1gm4s3pqlyvqqwaue44lqev2f9dx73d63ujj4zp":6,"iaa1gzhuyac7mt49pdccjyt70x85nfczz3t4fjxg57":6,"iaa1h2uh4zk348f2jg0jpcysj83n6h0htel0q8dads":6,"iaa1h6fr09724lzf6l3usarlfm2uulxzhn52stk6et":6,"iaa1h7mkj037j392mw5yaa4vx8ezjdzd7j49f2n4d4":6,"iaa1h99ay7d9mhch7l6dnd7ktcgyhmz8axsvzljcnu":6,"iaa1hmljfxzwu9v44rleykl5fezxpczyghp62k5l9d":6,"iaa1hmnwx3am7zpzxdgtdtgxngz3s0qzverm5lfsxa":6,"iaa1hrem0n6pq25appf2za7hkcc6ny2akh57j767nq":6,"iaa1hsd8pldjj8t75v53uyxnqsnfdv0t5jtszpa3wv":6,"iaa1j3ggzmck7tuqkq5u5ttjnqmgl2zhecz9ly6g6n":6,"iaa1j9sfx5zjuhhxphtd4ss0mfyj3ujcl0s85t5m96":6,"iaa1ja3u8vvknr5qmqqadeza2l9537ggf4xmsx2fn7":6,"iaa1ja7dxlymjmlfy89ftcg450qzn5carxh0fcq8yt":6,"iaa1jfwcmn27ejlxkupg2vzzzy74ze47vr9fsag9f7":6,"iaa1jhfh2zmm26lphxvgnwk3d7z8n6z3697arzn9e4":6,"iaa1jk3gae0v6hx759rh0473649n63hwll468hhhhy":6,"iaa1jq592yms9yr5yj33c0hedzlkhlugzetmg5haer":6,"iaa1jtn5qvtpjylmmfc3dkstq38a8r8qpwxus4eerv":6,"iaa1jxjgtxpdeaz09zn7wa7pzxf2sz6n3gw33su0a0":6,"iaa1jzmk5r5vgc4urwztrxzac52nm2apr8ck0atv8l":6,"iaa1k0nw6tr88mtfmddtw4e0rq9z7zdq6ua0uhmndv":6,"iaa1k2qpl2adl7j04xvwm8v7muepye64zx7l7aulrj":6,"iaa1k50uw7893nuxh3rk6wslguetu5clvj273ux2wu":6,"iaa1k6gfa9ry8g3hxngynyug7zhqcd00aagnhhn5hl":6,"iaa1k7h2dr4g69x33z85zkj850wdtdpfvkxdc3lvn3":6,"iaa1k9z3kvy5hxpcqx9rgg0wq6vm7srruqwatt7pfd":6,"iaa1kfvl8k00qahwvwstaemhwgge93h3jwljz742ca":6,"iaa1kpjx234szmapzt2j2dkpgfp0v5xhu304ntuuz5":6,"iaa1krkwzghswyu05ajunmrauqzvcvueycar37lt03":6,"iaa1ks2yqq5c5qpe8rhde5xg0m2ttndv3s2f6jlvg2":6,"iaa1ks75zlygm4zchc9684tvn0wxffs4fwf0p936rz":6,"iaa1ky9uxts805apn42a8g47rc050fvv8ca57q2uge":6,"iaa1l3668se0cw32q073rtw83r2245xn6ns5qtgh3a":6,"iaa1l76ppkjy4279rpz3l289xvh6v7e742d5de9q65":6,"iaa1l85zj9qy2tdlte49gyshmglznt95rmw9z2ap9j":6,"iaa1lelrwcmhghpvw98myvq2ahezysn782gyz4r6hx":6,"iaa1lrfacrtagv9sg22d3fqdfraj085gv29prj9jk8":6,"iaa1lwwn77r64lnvv4qu2zmuvvtluhz9444jw4cyad":6,"iaa1ly895td03fd8myeuzxfgwylrqnzgw7r54g9u37":6,"iaa1lzqw266aga0pqlavsfdhgzfgh8z42su0rg7jay":6,"iaa1m2y20wtuwyhawzjg7g3u04g64hr53tt7fe5zcd":6,"iaa1m32ueeud8xr9d42pnpc8j928wc2mvqw8w8jcte":6,"iaa1m5rv58qzyc9ucdv06wrs703n2rrefh5zmawem0":6,"iaa1m8wlegz29d5lmgqddwqzs6n40v3s4jawfj9qep":6,"iaa1mfw5qxr0je8ynm5wx0xa68npruqpngqzgpssye":6,"iaa1mmu82s7lhywjrn8u2lfcqt0uf8u6su80ssmhy7":6,"iaa1mmxdvnjd545c35zf770l9wz6fk6fy05p42g6mg":6,"iaa1mxt6c02756d4xaayk2rsxun5ajhntucky5daq0":6,"iaa1myanp9c4hs6dyunwadqflqmktws8uc3uhfxhd5":6,"iaa1n4d652h5jl8xy7mmj5mag99uc46w6xu2fzskzg":6,"iaa1n8aq7syawuh2pxvj2na0uvk9xqykz7jpelwny3":6,"iaa1nemurqs0d63z4cl9q8jhr0pd7a9md5vaxkwee2":6,"iaa1nfhjxg3kdpejynqmw6xdjfhdlqlqevh7mfakp7":6,"iaa1nfr94sqzu334qkzu3ckr5l7n3vz47svzcca5cx":6,"iaa1nphj0gcaq33qm42wl3969dasvggtlx98zxp9sm":6,"iaa1ntu9hqrqkr4qccleqzedyg7h70a6kp6x94s4tc":6,"iaa1nwwhvxe93fdduvhmny6f5wd3g4eyae0dycgj6z":6,"iaa1nz3wrs9el99g89s83zzdl6crkxjazu98fkaslk":6,"iaa1p24p5qrj4ym5vy99mnadsu7wlm0t63uhr27tfa":6,"iaa1p42d5pmcwrfx8mupx22qrlh3snqrxegrj0djy7":6,"iaa1p6ajdjmp6usncw97gsqvxafxsha9jn3qd7a8y8":6,"iaa1p8cnuklwu7hxzf66chszq7p559psmuf3649yu7":6,"iaa1pva93fdy0ftcrzcsfgef4aztkudnp43q4z7z8u":6,"iaa1q3t07luhmpzsjg7rd70t4yhk857hcarux2c23h":6,"iaa1q76tntcnepcnxyflt24yt7rrag579jy24knkj2":6,"iaa1qa048uw5awcc5cf5f9ezc5wkq4y4c0qkj7ujgd":6,"iaa1qhg5gx0m9v2zcml5s09r78yr66wwdan6tdapzq":6,"iaa1qnls07wrmpfpw698989a63m6mjpchp56a2f3pg":6,"iaa1qpj4qq23xcr7lc9d7pv5237e4yx6p6n2uv6dxz":6,"iaa1r2fjlvlhzc94pvg8fd5nunhazfhpxqpejlvm32":6,"iaa1r8exsx9ukul485mnq33cd949239ktn0k3t6lnt":6,"iaa1r8qkh46wqtkh8scc8gpu3pc2vm9z7pj5ef8hq2":6,"iaa1rc3efmd6llltssvtw5ar2w8sxadqtskkk4cyua":6,"iaa1rshjtzscwc3tw68mc28cym6ldn97ju72mxphm9":6,"iaa1rw7q6pf4yh4390czcjm0tjpaydgk0ells357ap":6,"iaa1rw8wlzscpguc5pajpt4ccdryz9pw48azuw73ct":6,"iaa1rweme2ttxg4z2scmmuxjkaeafx5t23sfmgk88f":6,"iaa1s02k2v4cy9wn2v57024259r7w6nc796ae6waq2":6,"iaa1s4jw93hhpywvejpsvagt9klv3zwlqw3dkwq9mx":6,"iaa1sdjq0jctzud0kckjtrum6xcdscmpmkqpwdvt9j":6,"iaa1snketpc7gced0x2lf94cgd6flhe76hqeh8q89n":6,"iaa1sqnk4w7xyfdv59as8xvszlmjct8h0wtkt6z02y":6,"iaa1ss59dcj7h5d5spuqxgsk7552l2kjm8nvm2jmkp":6,"iaa1sswv49c3tyka3e9mkzf7cxrhlmdn39760fsu96":6,"iaa1sw3dyq4c8uszdtplj7en8apu5y4f6xsxmpxr0l":6,"iaa1t325nrhc7xe3h86kg6y9d38nlt0vew5k8du7w4":6,"iaa1t4qwfh999eh3hjx5fgmczqncsdh7hj9k04zter":6,"iaa1t6mek4zy50jeuwk6tkj9gkmj8vyyfstt8qk9vf":6,"iaa1t7ahn9md5whyj0gwu6fvxs3ydrsqdjrzjde0zc":6,"iaa1ta8tzvu6k0hmcll0ay9z5nnpu3le95552pek99":6,"iaa1ter55q8kmwf3lz3kv3rq3z3dtupgne6v5u97av":6,"iaa1thrln2twkfft493htx2xsqxp0u6t2yt9wyqr2w":6,"iaa1ttf8vp9qsyn9zp8cl5sytph2f2pjukgys38rvd":6,"iaa1tvnu8j8j2p4yv7d83a8lxye6k38zlx8remskur":6,"iaa1tzjuk5af83nj57zs0nwe9jcl8tglfgxmvduz0l":6,"iaa1u4307vr77rn3ytvfqvucazwfqfl7kl2zqcn5ch":6,"iaa1u8safvtvjkthpwjy9r2372zyextgdau7h6td4h":6,"iaa1ufffgw3kvnw5nqnwx56vjnntds3f7hk95075dt":6,"iaa1uqj9dgqle0gkdswndhkmcxke02tczqv00r4vn8":6,"iaa1uqvzxed4c4z5h3rwdahvpyfwwd60c7p556rmh9":6,"iaa1v3w52dzqx9kau2eyqgc08npctrw085cphf54vy":6,"iaa1v7elddmwqhk0vsqselvddyez35rhay93lefzhy":6,"iaa1vmdygegeqnxklf9v207y92lqm5wavw73k8j4rk":6,"iaa1vrkhtgqn3qcaclpl4ukc7qafeg7khyjghhjevv":6,"iaa1vsj4e3rug05w86zgdtnqvvpkfeujwmygcse9t2":6,"iaa1vw0lfkv6neafsmk0k3rxm3eyrgtm2hsyy9plh2":6,"iaa1vxhufng59v46m6engtrpk2mjq2vpqzjm0z323r":6,"iaa1w3uect53nm3qmjl85l73qjh0vxwrs460q3vp3f":6,"iaa1w6nccltjkuwjcd5qmqlax7dfjuarmkqz0sqqe3":6,"iaa1whx9p9ae3nt8f8v8vkmc54846rclnc8lkla3ll":6,"iaa1wqh8adhg9lqwthgqyqnwgaddg9xszkh4ca2sd8":6,"iaa1wumgczmqrtq5twg6qvyfql0wdsakfddudj00r4":6,"iaa1wxfvw95ad2elxfxz8ltx29uhcwg0e0gn4qkrg2":6,"iaa1x0xk4yg4k0tmj8hrne8t32wuw8z3k2cq08y8e6":6,"iaa1x2q0acgyakg2j9q7ta0ntz0ps6t9xmr59hzxmx":6,"iaa1x6hmde5x9hadug9dqrcun7tkypva8r5qugse9m":6,"iaa1x90f3av778c8rltevdfnj3w3rcn70kw9n98tr4":6,"iaa1x9gt6tv94jqpsqep24y3qmzfmkmwpjvs4r7wdr":6,"iaa1xdhw4jpewpazcx8dgnt4y5dvg9s08zw6q6f7ee":6,"iaa1xkfhvaze0jdqtz8mw6a9dhgp630940r4x5qj3z":6,"iaa1xu5rw9dvdf7wdxmwjym398azxl57lrpsxazajt":6,"iaa1xzpcawt3m83hzpr377r5w2y3fwskqfh0fzd9rx":6,"iaa1y3ece7pp2xavrrjwwtezkg3zqnk5cdzzr327yt":6,"iaa1y9pvtw59w7636l3rymh36cjyey3mk968k2m37m":6,"iaa1yd23vlp9kgygq2wmuqkra0wgwa5jry3ffhw5ky":6,"iaa1yfwk32nuuxk4dk3707rk3tqz4rgaemnwvkdca5":6,"iaa1yge0mdrhc5xdsup4dfm00e3mhz4s9ck9klnszj":6,"iaa1yh3467zqjynemk46hf6329j0mhfus9dr4nwfa6":6,"iaa1ynnmgcyp72gxrw8nfdlw6e6mdnkanyunl8dkjw":6,"iaa1ynu2nv29zm8qge7kdecqjg5434jywdnq6lwajc":6,"iaa1ypeyntz6604pcmhzrtszgv8w0amxlczhq5kvr0":6,"iaa1ypy36vpw5fh2xffv7p2gj56ldhuqx30cgerqwc":6,"iaa1yqaamy8dwqehfmfdeedmc6j50sv5ymvx7lr9hc":6,"iaa1ysjhpc8aul7ej3rcvecqs2j2qy662f9w7es8au":6,"iaa1yst2uystk3s0j4zdlgvmkumlmw3grrutnk5ugr":6,"iaa1yts2l6lcjhycsq5a4k7aq7dcdgjp05td5ewllq":6,"iaa1ytxt3g79t350rs4ghysvjtj46l7m0s8n40032u":6,"iaa1ze3usqux85lk747kvaatj6493ywzsle9knjh69":6,"iaa1ze9qe7fw2ctd8d9ja8ag0gctuk4uvmkqtgy32a":6,"iaa1zfuvwe5xy6pf58wgt8dftjtn7jdwr9a42pwwr9":6,"iaa1zg46r7a08wl0k9yjdcx02camf7xmf7g3dns3kt":6,"iaa1zj5jq8gu2s4x8hr0l3yjq4hlsgyakwr4vstdam":6,"iaa1zl6vxewz4ev0ja2xdpv8dsdqfwl0hzpcdnchfq":6,"iaa1zld4lzspyvupfadfk0umlhh4ampmqw5h5yzy7d":6,"iaa1zu6w9m4vuy9gxkgg6wz8kvwr6l07ytslq79r5k":6,"iaa1zv6gt3syxahsk950p9xx0ef7xqlhsak70vl9nn":6,"iaa1zwtnc3k08n9gqxhwtxwzjdrkwgga8jp7qcx70m":6,"iaa1zwyzz99q8wwu7xjdxjhf3y8asmrzl8qxmr35pm":6}`
	if err := json.Unmarshal([]byte(jsonDataStr), &addrAmountMap); err != nil {
		t.Fatalf("unmarshal json fail, err is %s\n", err.Error())
	}

	handledDataStr := `["iaa1qtm7hw464wkr473k8qu30hgvh22gatdkc3w7p6","iaa1j7xaagn3ummuxwt9s9lwdykdk48clyvy85gcmr","iaa13tswxx3vzfg7apfh2l0jl54z87jlz6pq7ppw50","iaa1c6eh5zca9cv0495a8juas47hsyp5whn89zy529","iaa1y9pvtw59w7636l3rymh36cjyey3mk968k2m37m","iaa13248jes6rrswdvz8lljsda7dfelynsyjqrfcz0","iaa1d8k6qsgjt4ywytwj0qv52pg8h7y8j6uujankvs","iaa1rc3efmd6llltssvtw5ar2w8sxadqtskkk4cyua","iaa1s02k2v4cy9wn2v57024259r7w6nc796ae6waq2","iaa1vmdygegeqnxklf9v207y92lqm5wavw73k8j4rk","iaa1vxhufng59v46m6engtrpk2mjq2vpqzjm0z323r","iaa19tk96rvgqxe604m6zypyl7fuupkp4pckyayaeh","iaa1ja3u8vvknr5qmqqadeza2l9537ggf4xmsx2fn7","iaa1s4jw93hhpywvejpsvagt9klv3zwlqw3dkwq9mx","iaa1sw3dyq4c8uszdtplj7en8apu5y4f6xsxmpxr0l","iaa1y3ece7pp2xavrrjwwtezkg3zqnk5cdzzr327yt","iaa188q82c2v6ehc88r8ek0fn9mf3u5z92lcst9ulk","iaa15u9lkx734zcy6g5v8kkzr439ugar7ja2s0r86d","iaa18dvz99edrp3x6t4maxmyjnu02eg7r28zha5ur6","iaa14ha5urqqlk9phh4wmzjh5v9zldy663r8ys87ts","iaa1epq0wzld73c2vt3xngrak973mgth8fuazcw580","iaa12l9h8jmrymtdchejns3huegnern08rrel3kfr9","iaa19cmhg32kkc668sjsnakuvga5kq4d02prf5f6yu","iaa1rshjtzscwc3tw68mc28cym6ldn97ju72mxphm9","iaa19ch2jyl5v4jgazeprfv6v7jj6d70d6qmj6ex55","iaa19vwdxqkwrtnd4h7z5e9zz3qwa06jd0wn5n3qf9","iaa1gg6auuyv4h4pxfwsv2rj8p4myt0xj5a6l9x2ml","iaa1u4307vr77rn3ytvfqvucazwfqfl7kl2zqcn5ch","iaa132yd5q8vdmrnjz3ksdnmqfpx9q0zn8weq73edw","iaa1vrkhtgqn3qcaclpl4ukc7qafeg7khyjghhjevv","iaa1h99ay7d9mhch7l6dnd7ktcgyhmz8axsvzljcnu","iaa1ks75zlygm4zchc9684tvn0wxffs4fwf0p936rz","iaa1l76ppkjy4279rpz3l289xvh6v7e742d5de9q65","iaa1zu6w9m4vuy9gxkgg6wz8kvwr6l07ytslq79r5k","iaa1fchk7ahkmavsr7jscrua93w65rr5t4qyhuzduz","iaa1t7ahn9md5whyj0gwu6fvxs3ydrsqdjrzjde0zc","iaa1xkfhvaze0jdqtz8mw6a9dhgp630940r4x5qj3z","iaa1l85zj9qy2tdlte49gyshmglznt95rmw9z2ap9j","iaa1p42d5pmcwrfx8mupx22qrlh3snqrxegrj0djy7","iaa1q76tntcnepcnxyflt24yt7rrag579jy24knkj2","iaa1zld4lzspyvupfadfk0umlhh4ampmqw5h5yzy7d","iaa143suv89arrf6786gyu23446cxq7eysege7522a","iaa1jq592yms9yr5yj33c0hedzlkhlugzetmg5haer","iaa1p8cnuklwu7hxzf66chszq7p559psmuf3649yu7","iaa1rw7q6pf4yh4390czcjm0tjpaydgk0ells357ap","iaa1snketpc7gced0x2lf94cgd6flhe76hqeh8q89n","iaa1ss59dcj7h5d5spuqxgsk7552l2kjm8nvm2jmkp","iaa1aepg660fyg7lzwzkk2388zj80k69e4wdwn43qv","iaa1fke62y5ds9hkmnqkc0409qkz3dl4vj4jks2vrl","iaa1jtn5qvtpjylmmfc3dkstq38a8r8qpwxus4eerv","iaa1kfvl8k00qahwvwstaemhwgge93h3jwljz742ca","iaa1r8exsx9ukul485mnq33cd949239ktn0k3t6lnt","iaa1cdprulu3m5hk99el8c9wc4d208srrm627wcdr4","iaa1ky9uxts805apn42a8g47rc050fvv8ca57q2uge","iaa1wxfvw95ad2elxfxz8ltx29uhcwg0e0gn4qkrg2","iaa1654xshv4zyxarhtafw3483nfw04j8vs9g7pc9f","iaa17qhy6cuqyz4g4xefwaftzwgdse4va3z8y65fqa","iaa1qpj4qq23xcr7lc9d7pv5237e4yx6p6n2uv6dxz","iaa13e4rxrqgcgjwrnyeq99cz6run59tc794yxwrq4","iaa1ddckxhypgpdz0hnckd4evphqe0sqdscjnv0jxn","iaa1mmu82s7lhywjrn8u2lfcqt0uf8u6su80ssmhy7","iaa1ypy36vpw5fh2xffv7p2gj56ldhuqx30cgerqwc","iaa1lrfacrtagv9sg22d3fqdfraj085gv29prj9jk8","iaa1t325nrhc7xe3h86kg6y9d38nlt0vew5k8du7w4","iaa10cp5qa7hn6fj2xtaz7530n2jzfa7qmnsg2qdf9","iaa10qthtarg3mcgpezvr74gqaxv8hdms53mcw9a4x","iaa1j3ggzmck7tuqkq5u5ttjnqmgl2zhecz9ly6g6n","iaa1484ar040gd7ppq0rf667u7ptx8n0u6z4aswjkd","iaa1h7mkj037j392mw5yaa4vx8ezjdzd7j49f2n4d4","iaa1jhfh2zmm26lphxvgnwk3d7z8n6z3697arzn9e4","iaa1jk3gae0v6hx759rh0473649n63hwll468hhhhy","iaa1nemurqs0d63z4cl9q8jhr0pd7a9md5vaxkwee2","iaa1qa048uw5awcc5cf5f9ezc5wkq4y4c0qkj7ujgd","iaa17efdjm76myzfxw479eqwxsh9csc2565flgls2k","iaa1g3hsvqmp72lkuru33m4d9k292jlu7329e7fefr","iaa1pva93fdy0ftcrzcsfgef4aztkudnp43q4z7z8u","iaa1cnchyft0juwjkrj4cch9mxjdrq9zp8z9mtulmc","iaa1m2y20wtuwyhawzjg7g3u04g64hr53tt7fe5zcd","iaa1nz3wrs9el99g89s83zzdl6crkxjazu98fkaslk","iaa1a9g9hcwven940qwv8z2v5cpl3sneu9pzvze82y","iaa1r8qkh46wqtkh8scc8gpu3pc2vm9z7pj5ef8hq2","iaa1ttf8vp9qsyn9zp8cl5sytph2f2pjukgys38rvd","iaa1whx9p9ae3nt8f8v8vkmc54846rclnc8lkla3ll","iaa15gqfv3nmy4hqhj3hk0hlgdjlmejstewkygzwsz","iaa1cwypr8kslyr6naevr9s2djw85qk2dvsgrxkzdw","iaa1jxjgtxpdeaz09zn7wa7pzxf2sz6n3gw33su0a0","iaa1vw0lfkv6neafsmk0k3rxm3eyrgtm2hsyy9plh2","iaa1n4d652h5jl8xy7mmj5mag99uc46w6xu2fzskzg","iaa18fu2jj3ns95kjhjtz85gwzskqsr4hz47z904t2","iaa1rweme2ttxg4z2scmmuxjkaeafx5t23sfmgk88f","iaa12s0ddw23glasqh8e3pjpx40vcy3q9fhuqfq9xl","iaa12ufzlf78ptu8s9sn4lrwvertesefyhvvnacc65","iaa1yfwk32nuuxk4dk3707rk3tqz4rgaemnwvkdca5","iaa1d8ecgh2j9nrrzashs6geeftgq36qd7hdsnrql5","iaa1ks2yqq5c5qpe8rhde5xg0m2ttndv3s2f6jlvg2","iaa1mfw5qxr0je8ynm5wx0xa68npruqpngqzgpssye","iaa18skv95m9qqjklhtzzekdck4wht8ejt5k8qply4","iaa1a3dz7prw9q8jr54czscay26wnuw2knsfp6u0e9","iaa1cmd640lg0wu96s05rthzjt38agntpg5wz7x958","iaa1u8safvtvjkthpwjy9r2372zyextgdau7h6td4h","iaa1ysjhpc8aul7ej3rcvecqs2j2qy662f9w7es8au","iaa138kla03j07xj2rln5p3v5xvjrx8vvyk62y0rmw","iaa19d3kddcyz5ekwj9xnzg85met8damr80fgnelyh","iaa1lelrwcmhghpvw98myvq2ahezysn782gyz4r6hx","iaa1vsj4e3rug05w86zgdtnqvvpkfeujwmygcse9t2","iaa10padu9n2q7xqhdkde3umjtwpdrj9f3g5r4pjl6","iaa1k0nw6tr88mtfmddtw4e0rq9z7zdq6ua0uhmndv","iaa10tp6zc4lvcku5cfzltxl7y7vtwmuj9htg8yy5r","iaa1dkg9j2mmlmwg7uslaxc0hvn8e30grkedcte43s","iaa1nwwhvxe93fdduvhmny6f5wd3g4eyae0dycgj6z","iaa14nhqjude62xjsmv0wa7cagpwrvp293ddsuaeec","iaa18yr3jfdrkz58x6xpus9rxq35x8mq0wjp5w6trh","iaa1338yeawvta4a4zgtysez068u56met2k3pckp2u","iaa1j9sfx5zjuhhxphtd4ss0mfyj3ujcl0s85t5m96","iaa1mxt6c02756d4xaayk2rsxun5ajhntucky5daq0","iaa10mxldppuzcxjt23c07mtzx2xz3whc56ycrej2l","iaa17vyu0ewjl0kkzuczv795wck94julpw3xhl2l4v","iaa19tuuj9vj9n9jlrjyn4tp65ay7w5vptsnfwwdlk","iaa1t6mek4zy50jeuwk6tkj9gkmj8vyyfstt8qk9vf","iaa18n37ez6qp6sp6asacpw3cqw3sgv9033m2r7anh","iaa1r2fjlvlhzc94pvg8fd5nunhazfhpxqpejlvm32","iaa10gfg38wyn9r83uzrh4z3fytfz3qsxdtpdcq2mh","iaa14xtmx49cl7l92eqr5tz54taj8slzte25r2pqg9","iaa1n8aq7syawuh2pxvj2na0uvk9xqykz7jpelwny3","iaa1330008yykgwv07wk8te3a3wvjea6l8ypqddfxl","iaa18pxx4297wenyu7awgfy3nwghzgsvk5pc58hl35","iaa1e5zrx4yus4nkwy2v4utrz3uvrngan0pjmnhzu4","iaa1hmljfxzwu9v44rleykl5fezxpczyghp62k5l9d","iaa1k2qpl2adl7j04xvwm8v7muepye64zx7l7aulrj","iaa1kpjx234szmapzt2j2dkpgfp0v5xhu304ntuuz5","iaa1xzpcawt3m83hzpr377r5w2y3fwskqfh0fzd9rx","iaa167ug37cd6dmm0s0jdgjwpe0qwjhenh5u9za60k","iaa17v9q69ax8jftlgwuxn5mqakjz2jf230qkseqtf","iaa14u945vxnjhh8umnq8n3kvtay3em24dzmj05nvc","iaa1exrtdfy8czgfm8k2cuajdcyyfzecaecm99saqk","iaa1hrem0n6pq25appf2za7hkcc6ny2akh57j767nq","iaa1gm4s3pqlyvqqwaue44lqev2f9dx73d63ujj4zp","iaa1h2uh4zk348f2jg0jpcysj83n6h0htel0q8dads","iaa1sqnk4w7xyfdv59as8xvszlmjct8h0wtkt6z02y","iaa1ytxt3g79t350rs4ghysvjtj46l7m0s8n40032u","iaa1c9ptf0gewl9xjvrjw0ct9wtjk8q9mhcjpp4x2z","iaa1d0d3qf7kc8jct4zvsfuvpallanc5vkpa2ms0m5","iaa1eletxstz8aup6p8uvq0ku8lkg5999hzve0v3vf","iaa1026qrwaaejjavscnvepjm5t670n2s0ucffhed0","iaa1uqj9dgqle0gkdswndhkmcxke02tczqv00r4vn8","iaa13kk0xtedkunyl2a3u46p6um4ae0dz2rjatcye9","iaa1cmqgy38y62xplcsa4cqcucatretgdt66j0v9h4","iaa1zl6vxewz4ev0ja2xdpv8dsdqfwl0hzpcdnchfq","iaa10x0mva68vvqesgu4hd94ayfml0a382y53fl9hd","iaa124y3raq3flwt9n8e3qwjate9n2u2mrmep64yk0","iaa1c5ckcxlqdqrzwl0932kcmx82z90fu3q8gndzg3","iaa1zg46r7a08wl0k9yjdcx02camf7xmf7g3dns3kt","iaa1k50uw7893nuxh3rk6wslguetu5clvj273ux2wu","iaa1k7h2dr4g69x33z85zkj850wdtdpfvkxdc3lvn3","iaa1ypeyntz6604pcmhzrtszgv8w0amxlczhq5kvr0","iaa1p6ajdjmp6usncw97gsqvxafxsha9jn3qd7a8y8","iaa17dftdlzrzt9e09hmzr5cfs2n9zth52gyu080n0","iaa1ff0nyptye6esfx9tzk539knkkxxtcape86vwdt","iaa1m32ueeud8xr9d42pnpc8j928wc2mvqw8w8jcte","iaa1ntu9hqrqkr4qccleqzedyg7h70a6kp6x94s4tc","iaa1jzmk5r5vgc4urwztrxzac52nm2apr8ck0atv8l","iaa1yst2uystk3s0j4zdlgvmkumlmw3grrutnk5ugr","iaa1zwyzz99q8wwu7xjdxjhf3y8asmrzl8qxmr35pm","iaa19rahjgjq9n962l3ezkpa3vm6g898hatwax3kap","iaa19yjtks2wwrn0fqusku68es9l336e5z3gl62lne","iaa1jfwcmn27ejlxkupg2vzzzy74ze47vr9fsag9f7","iaa19atq5qj3snp000lz0mn40d7d2tkp350pkuu39m","iaa1hsd8pldjj8t75v53uyxnqsnfdv0t5jtszpa3wv","iaa1sdjq0jctzud0kckjtrum6xcdscmpmkqpwdvt9j","iaa180qcntxnh6lj0a6jmxfcpjf32k97u6l9xs8nnw","iaa1e9l6fhwl7jt0sddl63smprt34xsf5lszdc2g6n","iaa199xpvkx7sanrgqtufjetqjd9cknpjq9af4lrye","iaa1cmxh0zjqy9uq8c0vvd27z4x2cg7aajp7scrkh7","iaa1rw8wlzscpguc5pajpt4ccdryz9pw48azuw73ct","iaa1ze3usqux85lk747kvaatj6493ywzsle9knjh69","iaa17l0pwy4qvplz5s0zjg2tgp995999kna730fmts","iaa1m5rv58qzyc9ucdv06wrs703n2rrefh5zmawem0","iaa15h7awla6rem8ynj85p5wx8qujpcztmvppday3q","iaa18xgeau6ehq2qg2khh72egg9hq6xgtvvt7wmre0","iaa1h6fr09724lzf6l3usarlfm2uulxzhn52stk6et","iaa1lzqw266aga0pqlavsfdhgzfgh8z42su0rg7jay","iaa1ter55q8kmwf3lz3kv3rq3z3dtupgne6v5u97av","iaa1yd23vlp9kgygq2wmuqkra0wgwa5jry3ffhw5ky","iaa16qf9pn6u8kp2krzmw9g4qvw7l5ykamps0dezve","iaa1gzhuyac7mt49pdccjyt70x85nfczz3t4fjxg57","iaa1lwwn77r64lnvv4qu2zmuvvtluhz9444jw4cyad","iaa1v3w52dzqx9kau2eyqgc08npctrw085cphf54vy","iaa1wqh8adhg9lqwthgqyqnwgaddg9xszkh4ca2sd8","iaa10l0khu6ddhwlej94eptedznp7mj40g2enpe0s8","iaa10sp709q23q6kfgwcpc8syvf9culywr2y3wg2c4","iaa1cp04ymk4f08283fe4dqaaxfljvuurpnlqgr3ka","iaa1zv6gt3syxahsk950p9xx0ef7xqlhsak70vl9nn","iaa188zwzfvjqhc7fakyj7yw8d4fq24agmv0yadvzm","iaa1ft747dwuc9hn6xltmqmpmt4tdqs2ck37qgpjjs","iaa1t4qwfh999eh3hjx5fgmczqncsdh7hj9k04zter","iaa193t9ejjr63n0djw5wdkkvl0vqkuafew053ca90","iaa1yqaamy8dwqehfmfdeedmc6j50sv5ymvx7lr9hc","iaa1qhg5gx0m9v2zcml5s09r78yr66wwdan6tdapzq","iaa1w6nccltjkuwjcd5qmqlax7dfjuarmkqz0sqqe3","iaa1zwtnc3k08n9gqxhwtxwzjdrkwgga8jp7qcx70m","iaa1l3668se0cw32q073rtw83r2245xn6ns5qtgh3a","iaa1mmxdvnjd545c35zf770l9wz6fk6fy05p42g6mg","iaa102q4gpvh4kmaus402lgcn4u2ylfx80wr6e3ryt","iaa1ja7dxlymjmlfy89ftcg450qzn5carxh0fcq8yt","iaa1nfhjxg3kdpejynqmw6xdjfhdlqlqevh7mfakp7","iaa1q3t07luhmpzsjg7rd70t4yhk857hcarux2c23h","iaa1tzjuk5af83nj57zs0nwe9jcl8tglfgxmvduz0l","iaa1f93zjgn6dc6gl2yqjsh63cu6ykup63399gmy7c","iaa1307r6skc0d6qkl6r50kgxxeases0dh2zh0tllw","iaa1cyd8pmrdwllqrnavlklvr48qpm38zndnrn5gsa","iaa1nfr94sqzu334qkzu3ckr5l7n3vz47svzcca5cx","iaa13quvlsugg6jcsur55nh9njt82luvsnpg94ek5j","iaa1ly895td03fd8myeuzxfgwylrqnzgw7r54g9u37","iaa1yts2l6lcjhycsq5a4k7aq7dcdgjp05td5ewllq","iaa1zj5jq8gu2s4x8hr0l3yjq4hlsgyakwr4vstdam","iaa1xu5rw9dvdf7wdxmwjym398azxl57lrpsxazajt","iaa1fga5rw2mcr4wlduezlqk4ypmx53fqy32q2ngm0","iaa1v7elddmwqhk0vsqselvddyez35rhay93lefzhy","iaa1xdhw4jpewpazcx8dgnt4y5dvg9s08zw6q6f7ee","iaa12uu9j4v2p3clec429c9dxnyjk7z8v62cy5wtg5","iaa185vfkxw6cl6samd4thnw2g0kk7ueh0q4rhn27r","iaa1hmnwx3am7zpzxdgtdtgxngz3s0qzverm5lfsxa","iaa12fd49g2c8vf5gv8ljz6aut397knnnjmjgk8axe","iaa12u3g03ty5k6tplwek0zqknvcg6qtmslku26u05","iaa18v2z7er8mls0ceceu82jwefnjdvnw4xkmwlv75","iaa14ltzaah25djm49g4hqfcg3zujad8dzy66yfd6s","iaa1m8wlegz29d5lmgqddwqzs6n40v3s4jawfj9qep","iaa1uqvzxed4c4z5h3rwdahvpyfwwd60c7p556rmh9","iaa13dur05xudrpnmgukuca6zwluyn70kx39qcqg52","iaa140qzvm0g6v4et8wys53rcyy69huddh6savd5jk","iaa1k6gfa9ry8g3hxngynyug7zhqcd00aagnhhn5hl","iaa1x6hmde5x9hadug9dqrcun7tkypva8r5qugse9m","iaa10v9kf6u70rtqzew9l0a58j0xvq76r2zucmj78q","iaa12psd3exwxs2zgy9qrrmc37zn8pv2au8yjhs5lz","iaa1c88fm7faeqn59fvjwafuavdrmcwakv73fcyupf","iaa1x2q0acgyakg2j9q7ta0ntz0ps6t9xmr59hzxmx","iaa1cp9euwryr6mqlkufvxlmafmdatd0kwa4vgct5z","iaa1p24p5qrj4ym5vy99mnadsu7wlm0t63uhr27tfa","iaa1ufffgw3kvnw5nqnwx56vjnntds3f7hk95075dt","iaa1yge0mdrhc5xdsup4dfm00e3mhz4s9ck9klnszj","iaa12c8ckm9su3jamy9czqe7ysyfl49qk5p8af5dah","iaa1g22efp6lln6jqw5yhhjy2q5kdjhuphx46rdygy","iaa1myanp9c4hs6dyunwadqflqmktws8uc3uhfxhd5","iaa1ta8tzvu6k0hmcll0ay9z5nnpu3le95552pek99","iaa1tvnu8j8j2p4yv7d83a8lxye6k38zlx8remskur","iaa1ynu2nv29zm8qge7kdecqjg5434jywdnq6lwajc","iaa1x90f3av778c8rltevdfnj3w3rcn70kw9n98tr4","iaa1ze9qe7fw2ctd8d9ja8ag0gctuk4uvmkqtgy32a","iaa1x0xk4yg4k0tmj8hrne8t32wuw8z3k2cq08y8e6","iaa1x9gt6tv94jqpsqep24y3qmzfmkmwpjvs4r7wdr","iaa1glyjp4muty76g0rms9c5dx6hwhfmvtgpld2dvw","iaa1wumgczmqrtq5twg6qvyfql0wdsakfddudj00r4","iaa1em29a6gvud976nccc4gj8u7txr04z24t3rlg5e","iaa1nphj0gcaq33qm42wl3969dasvggtlx98zxp9sm","iaa19qdhqcu56x09jjs6e79pp35temq0hm0h4hfn6r","iaa1sswv49c3tyka3e9mkzf7cxrhlmdn39760fsu96","iaa16hvxxftn9lkt5fwdsamkn7fksy5d6zg8f4w5mr","iaa16xf95cpuk979p0d7tjterm09jgglscrx0xrjl3","iaa1ynnmgcyp72gxrw8nfdlw6e6mdnkanyunl8dkjw"]`
	var handledAddrList []string
	if err := json.Unmarshal([]byte(handledDataStr), &handledAddrList); err != nil {
		t.Fatalf("unmarshal json fail, err is %s\n", err.Error())
	}

	if len(addrAmountMap) == 0 {
		fmt.Println("data is empty")
		return
	}

	listToMap := func(list []string) map[string]bool {
		res := make(map[string]bool)
		if len(list) > 0 {
			for _, v := range list {
				res[v] = true
			}
		}
		return res
	}
	handledAddrMap := listToMap(handledAddrList)

	for k, v := range addrAmountMap {
		if handledAddrMap[k] {
			fmt.Printf("%s has been handled\n", k)
			continue
		}
		amount := util.Float64ToStr(v * math.Pow10(18))
		coin := ctypes.Coin{
			Denom:  "iris-atto",
			Amount: amount,
		}
		coins := []ctypes.Coin{coin}
		if res, err := txClient.SendToken(k, coins, memo, true); err != nil {
			fmt.Printf("%s fail, err: %s\n", k, err.Error())
			time.Sleep(time.Duration(5) * time.Second)
		} else {
			if res.CommitResult.CheckTx.IsErr() || res.CommitResult.DeliverTx.IsErr() {
				fmt.Printf("%s fail, res: %s\n", k, util.ToJsonIgnoreErr(res))
				time.Sleep(time.Duration(5) * time.Second)
			} else {
				fmt.Printf("%s:%v success, txHash: %s\n", k, v, res.CommitResult.Hash)
				fmt.Println("now sleep 10s")
				time.Sleep(time.Duration(10) * time.Second)
			}
		}

	}
}
