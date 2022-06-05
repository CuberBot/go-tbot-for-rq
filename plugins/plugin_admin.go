package plugins

import (
	"context"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	. "github.com/2mf8/tbotGo/public"
	"github.com/2mf8/tbotGo/utils"
	"github.com/2mf8/go-pbbot-for-rq"
	"math/rand"
	"github.com/2mf8/go-pbbot-for-rq/proto_gen/onebot"
)

type Admin struct {
}

func (admin *Admin) Do(ctx *context.Context, bot *pbbot.Bot, event *onebot.GroupMessageEvent) (retval uint) {
	rawMsg := strings.TrimSpace(event.RawMessage)
	groupId := event.GroupId
	userId := event.Sender.UserId
	botId := bot.BotId
	var duration int = 0

	rand.Seed(time.Now().UnixNano())
	//success := rand.Intn(101)
	//delete := rand.Intn(101) + 200
	failure := rand.Intn(101) + 400

	s, b := Prefix(rawMsg, ".")
	if b == false {
		return utils.MESSAGE_IGNORE
	}
	reg1 := regexp.MustCompile("<at qq=\"")
	reg2 := regexp.MustCompile("\"/>")
	reg3 := regexp.MustCompile("  ")

	str1 := strings.TrimSpace(reg1.ReplaceAllString(s, ""))
	str2 := strings.TrimSpace(reg2.ReplaceAllString(str1, " "))

	for Contains(str2, "  ") {
		str2 = strings.TrimSpace(reg3.ReplaceAllString(str2, " "))
	}

	if StartsWith(str2, "jin") && (IsAdmin(bot, groupId, userId) || IsBotAdmin(userId)) {
		str2 = strings.TrimSpace(string([]byte(str2)[len("jin"):]))
		str3 := strings.Split(str2, " ")

		if len(str3) != 2 {
			replyText := strconv.Itoa(failure) + "（禁言格式错误）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return utils.MESSAGE_BLOCK
		}
		jinId, err := strconv.ParseInt(str3[0], 10, 64)
		if err != nil {
			replyText := strconv.Itoa(failure) + "（禁言对象错误）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return utils.MESSAGE_BLOCK
		}

		reg4 := regexp.MustCompile("天")
		reg5 := regexp.MustCompile("小时")
		reg6 := regexp.MustCompile("时")
		reg7 := regexp.MustCompile("分")
		reg8 := regexp.MustCompile("秒")
		str4 := strings.TrimSpace(reg4.ReplaceAllString(str3[1], "d"))
		str4 = strings.TrimSpace(reg5.ReplaceAllString(str4, "h"))
		str4 = strings.TrimSpace(reg6.ReplaceAllString(str4, "h"))
		str4 = strings.TrimSpace(reg7.ReplaceAllString(str4, "m"))
		str4 = strings.TrimSpace(reg8.ReplaceAllString(str4, "s"))
		str4 = str4 + "s"

		reg9 := regexp.MustCompile(`([0-9]+)(d|h|m|s)`)
		m := reg9.FindAllString(str4, -1)
		for _, v := range m {
			if EndsWith(v, "d") {
				num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("d")]))
				duration += num * 60 * 60 * 24
			}
			if EndsWith(v, "h") {
				num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("h")]))
				duration += num * 60 * 60
			}
			if EndsWith(v, "m") {
				num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("m")]))
				duration += num * 60
			}
			if EndsWith(v, "s") {
				num, _ := strconv.Atoi(string([]byte(v)[:len(v)-len("s")]))
				duration += num
			}
		}
		if duration <= 0 {
			replyText := "解除 " + strconv.Itoa(int(jinId)) + " 的禁言"
			bot.SetGroupBan(groupId, jinId, 0)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		}
		d := int32(duration)
		if duration < 30*60*60*24 {
			replyText := "禁言 " + strconv.Itoa(int(jinId)) + " " + strconv.Itoa(int(d)) + "秒"
			bot.SetGroupBan(groupId, jinId, d)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		} else {
			replyText := strconv.Itoa(failure) + "禁言时间超过最大允许范围"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return utils.MESSAGE_BLOCK
		}
	}

	if (StartsWith(str2, "t") || StartsWith(str2, "T")) && (IsAdmin(bot, groupId, userId) || IsBotAdmin(userId)) {
		rejectAddAgain := StartsWith(str2, "T")
		str2 = strings.TrimSpace(string([]byte(strings.ToLower(str2))[len("t"):]))
		tId, err := strconv.ParseInt(str2, 10, 64)
		if err != nil {
			replyText := strconv.Itoa(int(failure)) + "（踢出对象错误）"
			replyMsg := pbbot.NewMsg().Text(replyText)
			log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
			_, _ = bot.SendGroupMessage(groupId, replyMsg, false)
			return utils.MESSAGE_BLOCK
		}
		replyText := "踢出 " + strconv.Itoa(int(tId)) + " 成功"
		bot.SetGroupKick(groupId, tId, rejectAddAgain)
		log.Printf("[INFO] Bot(%v) Group(%v) -> %v", botId, groupId, replyText)
		return utils.MESSAGE_BLOCK
	}
	return utils.MESSAGE_IGNORE
}

func init() {
	utils.Register("群管", &Admin{})
}