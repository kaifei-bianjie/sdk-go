package ad

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/irisnet/irishub-sdk-go/modules/nft"
	sdktypes "github.com/irisnet/irishub-sdk-go/types"
	"github.com/irisnet/sdk-go/util"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	nftNamePrefix              = "Bifrost Testnet Badge"
	nftStatus                  = "active"
	nftDataMedalFaceValueDenom = "uiris"
	nftDenom                   = "bifrostestnetbadge"
	nftDenomName               = "BifrostTestnetBadges"
	nftRedeemable              = true
	nftImgKeyPrefix            = "nft/badge/"
	nftVoucherDataFilePath     = "./nft_badges_modify_2.xlsx"
)

var (
	faceValueMap = map[string]nftDataMedalFaceValue{
		"gold": nftDataMedalFaceValue{
			Denom:  nftDataMedalFaceValueDenom,
			Amount: "300000000",
		},
		"silver": nftDataMedalFaceValue{
			Denom:  nftDataMedalFaceValueDenom,
			Amount: "201000000",
		},
		"bronze": nftDataMedalFaceValue{
			Denom:  nftDataMedalFaceValueDenom,
			Amount: "99000000",
		},
	}
	badgeLevelNameMap = map[string]int{
		"gold":   1,
		"silver": 2,
		"bronze": 3,
	}
	badgeIconMap = map[string]string{
		"gold":   "https://rb-app.oss-cn-shanghai.aliyuncs.com/icon_badge_gold.png",
		"silver": "https://rb-app.oss-cn-shanghai.aliyuncs.com/icon_badge_silver.png",
		"bronze": "https://rb-app.oss-cn-shanghai.aliyuncs.com/icon_badge_bronze.png",
	}
	nftImgSvgSourcePathMap = map[string]string{
		"gold":   "./source/svg_data_gold.txt",
		"silver": "./source/svg_data_silver.txt",
		"bronze": "./source/svg_data_bronze.txt",
	}
	nftImgSvgTemplatePath = "./source/svg_img_template.svg"
)

type (
	nftDataMedalFaceValue struct {
		Denom  string `json:"denom"`
		Amount string `json:"amount"`
	}
	nftDataMedal struct {
		Holder          string                `json:"holder"`
		Status          string                `json:"status"`
		Level           int                   `json:"level"`
		Img             string                `json:"img"`
		Icon            string                `json:"icon"`
		FaceValue       nftDataMedalFaceValue `json:"face_value"`
		RedeemStartTime int                   `json:"redeem_start_time"`
		RedeemEndTime   int                   `json:"redeem_end_time"`
		Redeemable      bool                  `json:"redeemable"`
	}

	nftData struct {
		Medal nftDataMedal `json:"medal"`
	}
)

func TestNftMint(t *testing.T) {
	initClient()
	initOssClient()
	data := nftData{}                    //nft data
	medal := nftDataMedal{}              //nft data medal property
	levelCounter := make(map[string]int) //计数
	memo := ""
	getAccInfo := func(addr string) (uint64, uint64, error) {
		var accNumber, sequence uint64
		if acc, err := irisClient.QueryAccount(addr); err != nil {
			err := fmt.Errorf("get sender acc info fail, addr:%s, err:%s", addr, err.Error())
			return accNumber, sequence, err
		} else {
			accNumber = acc.AccountNumber
			sequence = acc.Sequence
		}
		return accNumber, sequence, nil
	}
	var signerAccNumber, signerSequence uint64
	sender := sendAddr
	if v1, v2, err := getAccInfo(sender); err != nil {
		fmt.Println(err.Error())
		return
	} else {
		signerAccNumber = v1
		signerSequence = v2
	}

	baseTx := sdktypes.BaseTx{
		From:          fromName,
		Password:      fromPassword,
		Gas:           gasLimit,
		Fee:           sdktypes.DecCoins{fee},
		Memo:          memo,
		Mode:          sdktypes.Sync,
		AccountNumber: signerAccNumber,
		Sequence:      signerSequence,
	}

	request := nft.MintNFTRequest{}
	xlsx, err := excelize.OpenFile(nftVoucherDataFilePath)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// 获取 Sheet1 上所有单元格
	excel := xlsx.GetRows("Sheet1")
	header := excel[0]
	if len(header) < 4 {
		panic("excel 数据格式有错误")
	}

	index := 1
	var nftImgSvgTemplateStr string                      //nft image template
	nftImgSvgSourceContentMap := make(map[string]string) //图片内容
	for i := 1; i < len(excel); i++ {
		row := excel[i]
		holder := row[0]
		if holder == "" {
			continue
		}
		startTime, _ := strconv.Atoi(row[1])
		endTime, _ := strconv.Atoi(row[2])
		for j := 3; j < len(row); j++ {
			level := strings.Split(header[j], "_")[1]
			quantity, _ := strconv.Atoi(row[j])
			for k := 0; k < quantity; k++ {
				if index%100 == 0 {
					time.Sleep(10 * time.Second)
				}
				//整理数据，组装参数
				nftId, nftName, nftImgKey := genNftInfo(holder, levelCounter, level)
				nftDataMedalBuilder(&medal, holder, level, nftImgKey, startTime, endTime)
				data.Medal = medal
				jsonData, _ := json.Marshal(data)
				request.Denom = nftDenom
				request.URI = ""
				request.ID = nftId
				request.Name = nftName
				request.Data = string(jsonData)
				//mint nft
				if res, err := irisClient.NFT.MintNFT(request, baseTx); err != nil {
					t.Errorf("nft mint fail, index: %d, nftId:%s, nftInfo:%s, errCode: %d, codeSpace: %s, err: %s\n",
						index, request.ID, util.ToJsonIgnoreErr(request), err.Code(), err.Codespace(), err.Error())
					time.Sleep(time.Duration(5) * time.Second)
					if _, v2, err := getAccInfo(sender); err != nil {
						t.Errorf("query acc info fail, addr:%s, err:%s\n", sender, err.Error())
					} else {
						baseTx.Sequence = v2
					}
				} else {
					t.Logf("nft mint success, txHash: %s, nftId: %s\n", res.Hash, nftId)
					baseTx.Sequence += 1
					//发交易成功，上传图片
					uploadNftImg(&nftImgSvgTemplateStr, nftImgSvgSourceContentMap, nftImgKey, level, nftId)
				}
				index++
			}
		}
	}
	t.Log("TestNftMint finish")
	for s, v := range levelCounter {
		t.Logf("%s:%d", s, v)
	}
}

var (
	denomId      = flag.String("denomId", "", "")
	nftId        = flag.String("nftId", "", "")
	nftName      = flag.String("nftName", "", "")
	nftRecipient = flag.String("nftRecipient", "", "")
	uri          = flag.String("uri", "", "")
	data         = flag.String("data", "", "")
)

func TestNftMintSingle(t *testing.T) {
	initClient()
	baseTx := sdktypes.BaseTx{
		From:     fromName,
		Password: fromPassword,
		Gas:      gasLimit,
		Fee:      sdktypes.DecCoins{fee},
		Memo:     "test",
		Mode:     sdktypes.Sync,
	}
	flag.Parse()
	request := nft.MintNFTRequest{}
	request.Denom = *denomId
	request.URI = *uri
	request.ID = *nftId
	request.Name = *nftName
	request.Data = *data
	request.Recipient = *nftRecipient
	if res, err := irisClient.NFT.MintNFT(request, baseTx); err != nil {
		t.Errorf("nft mint fail, errCode: %d, codeSpace: %s, err: %s\n", err.Code(), err.Codespace(), err.Error())
	} else {
		t.Logf("nft mint success, txHash: %s", res.Hash)
	}
}

func TestNftEditSingle(t *testing.T) {
	initClient()
	baseTx := sdktypes.BaseTx{
		From:     fromName,
		Password: fromPassword,
		Gas:      gasLimit,
		Fee:      sdktypes.DecCoins{fee},
		Memo:     "test",
		Mode:     sdktypes.Sync,
	}
	flag.Parse()
	request := nft.EditNFTRequest{
		Denom: *denomId,
		URI:   *uri,
		ID:    *nftId,
		Name:  *nftName,
		Data:  *data,
	}
	if res, err := irisClient.NFT.EditNFT(request, baseTx); err != nil {
		t.Errorf("nft edit fail, errCode: %d, codeSpace: %s, err: %s\n", err.Code(), err.Codespace(), err.Error())
	} else {
		t.Logf("nft edit success, txHash: %s", res.Hash)
	}
}

func TestNftBurnSingle(t *testing.T) {
	initClient()
	baseTx := sdktypes.BaseTx{
		From:     fromName,
		Password: fromPassword,
		Gas:      gasLimit,
		Fee:      sdktypes.DecCoins{fee},
		Memo:     "test",
		Mode:     sdktypes.Sync,
	}
	flag.Parse()
	request := nft.BurnNFTRequest{
		Denom: *denomId,
		ID:    *nftId,
	}
	if res, err := irisClient.NFT.BurnNFT(request, baseTx); err != nil {
		t.Errorf("nft burn fail, errCode: %d, codeSpace: %s, err: %s\n", err.Code(), err.Codespace(), err.Error())
	} else {
		t.Logf("nft burn success, txHash: %s", res.Hash)
	}
}

func genNftInfo(holder string, levelCounter map[string]int, level string) (nftId, nftName, nftImgKey string) {
	counter, ok := levelCounter[level]
	if ok {
		counter += 1
		levelCounter[level] = counter
	} else {
		counter = 1
		levelCounter[level] = counter
	}
	counterItoa := strconv.Itoa(counter)
	for i := len(counterItoa); i < 4; i++ {
		counterItoa = "0" + counterItoa
	}
	nftId = fmt.Sprintf("%s%s%s", level, counterItoa, holder[len(holder)-8:])
	nftName = fmt.Sprintf("%s %s %s %s", nftNamePrefix, level, counterItoa, holder[len(holder)-8:])
	nftImgKey = fmt.Sprintf("%s%s%s.svg", nftImgKeyPrefix, nftDenom, nftId)
	return nftId, nftName, nftImgKey
}

func nftDataMedalBuilder(medal *nftDataMedal, holder, level, nftImgKey string, redeemStartTime, redeemEndTime int) {
	medal.Holder = holder
	medal.Level, _ = badgeLevelNameMap[level]
	medal.Img = genOssUrl(nftImgKey)
	medal.Icon = badgeIconMap[level]
	medal.FaceValue = faceValueMap[level]
	medal.Status = nftStatus
	medal.Redeemable = nftRedeemable
	medal.RedeemEndTime = redeemEndTime
	medal.RedeemStartTime = redeemStartTime
}

func uploadNftImg(nftImgSvgTemplateStr *string, nftImgSvgSourceContentMap map[string]string, imgKey, level, nftId string) {
	if *nftImgSvgTemplateStr == "" {
		if c, err := util.ReadFile(nftImgSvgTemplatePath); err != nil {
			fmt.Printf("读取文件失败，file path：%s, error: %s\n", nftImgSvgTemplatePath, err.Error())
			return
		} else {
			*nftImgSvgTemplateStr = c
		}
	}
	uploadStr := *nftImgSvgTemplateStr
	imgContent, ok := nftImgSvgSourceContentMap[level]
	if !ok {
		if c, err := util.ReadFile(nftImgSvgSourcePathMap[level]); err != nil {
			fmt.Printf("读取文件失败，file path：%s, error: %s\n", nftImgSvgSourcePathMap[level], err.Error())
			return
		} else {
			nftImgSvgSourceContentMap[level] = c
			imgContent = c
		}
	}
	uploadStr = strings.Replace(uploadStr, "${img_data}", imgContent, 1)
	uploadStr = strings.Replace(uploadStr, "${nft_id}", nftId, 1)
	option := oss.ContentType("text/xml")
	objectAcl := oss.ObjectACL(oss.ACLPublicRead)

	err := ossBucket.PutObject(imgKey, strings.NewReader(uploadStr), option, objectAcl)
	if err != nil {
		fmt.Printf("upload img error,Error:%s\n", err)
		return
	} else {
		//fmt.Printf("上传图片成功,imgKey:%s\n", imgKey)
	}
}

func TestNftIssue(t *testing.T) {
	schemaJsonStr := `{"$schema":"http://json-schema.org/draft-07/schema#","$id":"http://json-schema.org/draft-07/schema#","title":"Medal Denom schema meta-schema","type":"object","properties":{"medal":{"type":"object","properties":{"holder":{"type":"string"},"img":{"type":"string","description":"voucher img of voucher holder"},"icon":{"type":"string","description":"voucher icon of voucher holder"},"status":{"type":"string","enum":["active","redeemed","expired"]},"level":{"type":"integer","description":"level of voucher, according to denom descript level meaning","enum":[1,2,3]},"face_value":{"type":"object","properties":{"denom":{"type":"string"},"amount":{"type":"string"}}},"redeemable":{"type":"boolean","enum":[true,false]},"redeem_after":{"type":"string"},"redeem_before":{"type":"string"},"redeem_start_time":{"type":"integer"},"redeem_end_time":{"type":"integer"}},"required":["holder","img","icon","status","level","face_value","redeemable","redeem_start_time","redeem_end_time"]}},"definitions":{}}`

	initClient()
	baseTx := sdktypes.BaseTx{
		From:     fromName,
		Password: fromPassword,
		Gas:      gasLimit,
		Fee:      sdktypes.DecCoins{fee},
		Memo:     "",
		Mode:     sdktypes.Sync,
	}
	request := nft.IssueDenomRequest{
		ID:     nftDenom,
		Name:   nftDenomName,
		Schema: schemaJsonStr,
	}
	if res, err := irisClient.NFT.IssueDenom(request, baseTx); err != nil {
		t.Fatalf("nft mint fail, errCode: %d, codeSpace: %s, err: %s", err.Code(), err.Codespace(), err.Error())
	} else {
		t.Logf("issue denom success, txHash: %s, denomId: %s, denomName: %s", res.Hash, request.ID, request.Name)
	}
}
