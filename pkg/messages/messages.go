package messages

// generic controls
const (
	Approve = "✅ OK"
	Reject  = "❌ Закрыть"
	Back    = "🔙 Назад"
)

// main menu
const (
	Cb5          = "😈 Добавить дроп с 5 КБ"
	Cb6          = "👹 Добавить дроп с 6 КБ"
	Cb4          = "👾 Добавить дроп с 4 КБ"
	Stats        = "📈 Моя статистика"
	Notification = "Оповещения"
	FullStats    = "📊 Статистика всех пользователей"
	Help         = "ℹ️ Помощь и обратная связь"
)

// add drop inline menu
const (
	Clear        = "🔄 Заново"
	AncientShard = "💙 Древний"
	VoidShard    = "💜 Темный"
	SacredShard  = "💛 Сакрал"
	EpicTome     = "📘 Эпик том"
	LegTome      = "📙 Лег том"
	Nothing      = "😭 Ничего"
)

// stats menu
const (
	LastVoidShard   = "💜 Последний войд"
	LastSacredShard = "💛 Последний сакрал"
	LastLegTome     = "📙 Последний лег том"
	MonthStats      = "📅 Весь мой дроп"
	MonthStatsShort = "🗓️ Суммарный дроп"
	FullStatExcel   = "💾 Выгрузка статистики"
)

// month menu
const (
	Jan = "Январь"
	Feb = "Февраль"
	Mar = "Март"
	Apr = "Апрель"
	May = "Май"
	Jun = "Июнь"
	Jul = "Июль"
	Aug = "Август"
	Sep = "Сентябрь"
	Oct = "Октябрь"
	Nov = "Ноябрь"
	Dec = "Декабрь"
)

const (
	Days30 = "30 дней"
	Days7  = "7 дней 📼📞📺👻☠️"
)

// help and faq
const (
	HelpHeader = "Привет!🤖\n" +
		"Если ты предпочитаешь видео-инструкции, то посмотри ролик обо мне на [канале LesiQ](https://www.youtube.com/watch?v=Va27Po7mmkU), а если текстовые – просто читай дальше.\n\n" +
		"Я бот для отслеживания твоего дропа с КБ в Raid SL.\n" +
		"Если ты забираешь *2 последних* сундука с *4, 5 и/или 6* клан босса, отправляй мне информацию о своем дропе, и я запомню ее для тебя. А еще покажу тебе разную интересную статистику.\n" +
		"Пока что *я умею*\n" +
		"  – показывать твой дроп за месяц,\n" +
		"  – показывать даты выпадения последних сакрала, войда и лег тома,\n" +
		"  – напоминать тебе о том, что нужно записать дроп (а еще побить кб), кстати, ты можешь сам(а) установить время оповещения или отключить его,\n" +
		"  – рассказывать тебе, насколько ты удачлив(а) по сравнению со всеми пользователями.\n" +
		"  – экспортировать твой дроп в Excel файл,\n" +
		"Но совсем *скоро я смогу*\n" +
		"  – показывать шансы дропа,\n" +
		"  – рассказывать тебе, насколько ты удачлив(а) по сравнению с сокланами.\n"
	HelpFAQ = "*Как добавить мой дроп?*\n" +
		"1. Нажми на кнопку _\"👾/😈/👹 Добавить дроп с 4/5/6 КБ\"_.\n" +
		"2. Нажимая на _кнопки осколков и томов_ (💙💜💛📘📙), набери свой дроп. Проверь, что все правильно в сообщении над кнопками. На кнопки можно нажимать несколько раз. Если что-то ввелось неверно, нажми на кнопку _\"🔄 Заново\"_.\n" +
		"3. Нажми на кнопку _\"✅ ОК\"_.\n" +
		"⚠️ Обрати внимание, что дроп нужно добавить до полуночи по МСК, иначе он запишется на следующий день.\n\n" +
		"*Добавил(а) что-то не то, что делать?*\n" +
		"– Если ты еще не нажал(а) на кнопку _\"✅ ОК\"_, то просто нажми на кнопку _\"🔄 Заново\"_ и снова введи дроп.\n" +
		"– Если ты сегодня уже отправил(а) информацию о дропе боту, то просто снова отправь дроп, и он перезапишет старый.\n" +
		"⚠️ Обрати внимание, что дроп можно обновить только до полуночи по МСК.\n\n" +
		"*Мне ничего не упало с 5/4 кб, что делать?*\n" +
		"_Записывать пустой дроп важно, это поможет мне правильно считать статистику._\n" +
		"1. Нажми на кнопку _\"👾/😈 Добавить дроп с 4/5 КБ\"_.\n" +
		"2. Нажми на кнопку _\"😭 Ничего\"_.\n\n" +
		"*Как записать дроп за вчера (неделю назад, 01.01.2001)?*\n" +
		"Никак. Дроп записывается в текущую дату, причем дата обновляется в полночь по МСК.\n\n" +
		"*Я бью 1/2/3 КБ, куда добавить мой дроп?*\n" +
		"Никуда. По моим расчетам, начинающие игроки уже спустя несколько месяцев начинают бить 4го и даже 5го КБ. Я уверен, что у тебя все получится, и что мы увидимся совсем скоро!\U0001F9BE\n\n" +
		"*Хочу записывать только сакралы (лег тома, ...), только один сундук, так можно?*\n" +
		"Да, так можно. Но совсем скоро я научусь считать разную глобальную статистику (шанс дропа, твою удачливость,...) с учетом того, что все записывают весь дроп из 2 сундуков. Чем больше людей делают иначе, тем менее точной получится статистика.\n" +
		"*Как установить время оповещения?*\n" +
		"Напиши боту /notification\\_on\\_fill\\_drop и дальше через пробел время по МСК в формате HH:MM. Например, /notification\\_on\\_fill\\_drop _13:30_ \n" +
		"*Как отключить оповещение?*\n" +
		"Напиши боту /notification\\_off\\_fill\\_drop и оповещения тебя больше не побеспокоят.\n\n"
)

// Full stats
const (
	FullStatsMsg = "Вот обощненная статистика по месяцам\n" +
		"Ноябрь 2021 https://telegra.ph/Statistika-po-KB-11-19-2"
)
