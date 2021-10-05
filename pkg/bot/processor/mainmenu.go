package processor

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	chatutils "vkokarev.com/rslbot/pkg/chat_utils"
	"vkokarev.com/rslbot/pkg/entities"
	"vkokarev.com/rslbot/pkg/keyboards"
	"vkokarev.com/rslbot/pkg/messages"
)

type MainProcessor struct{}

func (p *MainProcessor) Handle(ctx context.Context, state entities.UserState, msg *ProcessingMessage) (entities.UserState, []tgbotapi.Chattable, error) {
	switch msg.Text {
	case messages.Cb5:
		state.State = entities.StateCb5
		resp := chatutils.TextTo(msg, "Что упало с 5го КБ?", keyboards.AddDropInlineKeyboard)
		return state, resp, nil
	case messages.Cb6:
		state.State = entities.StateCb6
		resp := chatutils.TextTo(msg, "Что упало с 6го КБ?", keyboards.AddDropInlineKeyboard)
		return state, resp, nil
	case messages.Stats:
		state.State = entities.StateStats
		resp := chatutils.TextTo(msg, "Что тебе показать?", keyboards.StatsKeyboard)
		return state, resp, nil
	case messages.Help:
		text := "Привет!🤖\n" +
			"Если ты предпочитаешь видео-инструкции, то посмотри ролик обо мне на [канале LesiQ](https://www.youtube.com/watch?v=Va27Po7mmkU), а если текстовые – просто читай дальше.\n\n" +
			"Я бот для отслеживания твоего дропа с КБ в Raid SL.\n" +
			"Если ты забираешь *2 последних* сундука с *5 и/или 6* клан босса, отправляй мне информацию о своем дропе, и я запомню ее для тебя. А еще покажу тебе разную интересную статистику.\n" +
			"Пока что *я умею* показывать\n" +
			"  – твой дроп за месяц,\n" +
			"  – даты выпадения последних сакрала, войда и лег тома.\n" +
			"Но совсем *скоро я смогу*\n" +
			"  – экспортировать твой дроп в Excel файл,\n" +
			"  – показывать шансы дропа,\n" +
			"  – рассказывать тебе, насколько ты удачлив(а) по сравнению с сокланами и всеми пользователями.\n" +
			"Ты можешь прислать свои идеи по дальнейшему развитию [моему ботюне](https://t.me/rsl_cb_drop_support_bot).\n\n" +
			"А дальше я постараюсь ответить на *часто задаваемые вопросы*\n\n" +
			"*Как добавить мой дроп?*\n" +
			"1. Нажми на кнопку _\"😈/👹 Добавить дроп с 5/6 КБ\"_.\n" +
			"2. Нажимая на _кнопки осколков и томов_ (💙💜💛📘📙), набери свой дроп. Проверь, что все правильно в сообщении над кнопками. На кнопки можно нажимать несколько раз. Если что-то ввелось неверно, нажми на кнопку _\"🔄 Заново\"_.\n" +
			"3. Нажми на кнопку _\"✅ ОК\"_.\n" +
			"⚠️ Обрати внимание, что дроп нужно добавить до полуночи по МСК, иначе он запишется на следующий день.\n\n" +
			"*Добавил(а) что-то не то, что делать?*\n" +
			"– Если ты еще не нажал(а) на кнопку _\"✅ ОК\"_, то просто нажми на кнопку _\"🔄 Заново\"_ и снова введи дроп.\n" +
			"– Если ты сегодня уже отправил(а) информацию о боте, то просто снова отправь дроп, и он перезапишет старый.\n" +
			"⚠️ Обрати внимание, что дроп можно обновить только до полуночи по МСК.\n\n" +
			"*Мне ничего не упало с 5 кб, что делать?*\n" +
			"_Записывать пустой дроп важно, это поможет мне правильно считать статистику._" +
			"1. Нажми на кнопку _\"😈 Добавить дроп с 5 КБ\"_.\n" +
			"2. Нажми на кнопку _\"✅ ОК\"_.\n\n" +
			"*Как записать дроп за вчера (неделю назад, 01.01.200)?*\n" +
			"Никак. Дроп записывается в текущую дату, причем дата обновляется в полночь по МСК.\n\n" +
			"*Я бью 1/2/3/4 КБ, куда добавить мой дроп?*\n" +
			"Никуда. По моим расчетам, начинающие игроки уже спустя несколько месяцев начинают бить 5го КБ. Я уверен, что у тебя все получится, и что мы увидимся совсем скоро!\U0001F9BE\n\n" +
			"*Хочу записывать только сакралы (лег тома, ...), только один сундук, так можно?*\n" +
			"Да, так можно. Но совсем скоро я научусь считать разную глобальную статистику (шанс дропа, твою удачливость,...) с учетом того, что все записывают весь дроп из 2 сундуков. Чем больше людей делают иначе, тем менее точной получится статистика.\n\n" +
			"Если у тебя остались вопросы, напиши [моему ботюне](https://t.me/rsl_cb_drop_support_bot).\n"
		resp := chatutils.TextTo(msg, text, keyboards.MainMenuKeyboard)
		return state, resp, nil
	default:
		return state, nil, UnknownResuest
	}
}

func (p *MainProcessor) CanHandle(msg *ProcessingMessage) bool {
	switch msg.Text {
	case messages.Cb5, messages.Cb6, messages.Stats, messages.Help:
		return true
	default:
		return false
	}
}

func (p *MainProcessor) CancelFor(userID int64) {
}
