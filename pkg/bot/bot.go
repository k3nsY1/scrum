package bot

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gopkg.in/mgo.v2/bson"
	tb "gopkg.in/tucnak/telebot.v2"
)

var (
	menuMember = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnPlan    = menuMember.Text("Ввести данные для Scrum")

	menuData = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
	btnYes   = menuData.Text("Да")
	btnNo    = menuData.Text("Нет")
)

type Scrumdata struct {
	Department string
	Name       string
	Date       string
	Plan       string
}

//Структура бота
type Bot struct {
	Bot    *tb.Bot
	Member *tb.User
	Data   []Scrumdata
}

// Создание бота
func CreateBot(bot *tb.Bot) *Bot {
	return &Bot{
		Bot: bot,
	}
}

//Данные по сотруднику
func (b *Bot) collector() {
	r := Scrumdata{}
	b.Bot.Handle(&btnPlan, func(m *tb.Message) {
		b.Bot.Send(b.Member, "Введите свой департамент!")
		b.Bot.Handle(tb.OnText, func(m *tb.Message) {
			r.Department = m.Text
			b.Bot.Send(b.Member, "Введите свое имя")
			b.Bot.Handle(tb.OnText, func(m *tb.Message) {
				r.Name = m.Text
				b.Bot.Send(b.Member, "Введите сегодняшнюю дату")
				b.Bot.Handle(tb.OnText, func(m *tb.Message) {
					r.Date = m.Text
					b.Bot.Send(b.Member, "Введите задачи на сегодня")
					b.Bot.Handle(tb.OnText, func(m *tb.Message) {
						r.Plan = m.Text

						//Добавление в слайс
						b.Data = append(b.Data, r)
						b.endReport()
						b.Bot.Send(b.Member, "Правильно?", menuData)

					})
				})
			})
		})
	})
}

// Set client options
func (b *Bot) mkCol() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	scrumCollection := client.Database("db").Collection("scrum")

	// Insert One document
	scrum1 := b.Data[len(b.Data)-1]
	insertResult, err := scrumCollection.InsertOne(context.TODO(), scrum1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	cursor, err := scrumCollection.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var tbs []Scrumdata
	if err = cursor.All(context.TODO(), &tbs); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Found multiple documents: %+v\n", tbs)
}

// Get a handle for your collection
// collection := client.Database("scrum").Collection("scrum")

// //Подключение к БД
// client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
// if err != nil {
// 	log.Fatal(err)
// }
// ctx := context.Background()
// err = client.Connect(ctx)
// if err != nil {
// 	log.Fatal(err)
// }

// fmt.Println("Connected to MongoDB!")
// defer client.Disconnect(ctx)
// //Создание БД и Коллекции
// 	scrumDB := client.Database("scrum")
// 	err = scrumDB.CreateCollection(ctx, "scrum")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	scrumCollection := scrumDB.Collection("scrum")
// 	defer scrumCollection.Drop(ctx)
// 	result, err := scrumCollection.InsertOne(ctx, )
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("result", result)
// }

//Запись в json
// func (b *Bot) saveData() {
// 	rawDataOut, err := json.MarshalIndent(&b.Data, "", "  ")
// 	if err != nil {
// 		fmt.Println("ERROR")
// 		log.Fatal("JSON marshaling failed:", err)
// 	}
// 	fmt.Println(string(rawDataOut))
// 	err = ioutil.WriteFile("data.json", rawDataOut, fs.FileMode(os.O_APPEND))
// 	if err != nil {
// 		fmt.Println("ERROR")
// 		log.Fatal("Cannot write updated settings file:", err)
// 	}
// }

// Проверка данных
func (b *Bot) checkData() {
	b.Bot.Handle(&btnYes, func(m *tb.Message) {
		b.Bot.Send(b.Member, "Данные успешно записаны", menuMember)
		fmt.Println(b.Data)
		b.mkCol()
	})
	// b.saveData()
	b.Bot.Handle(&btnNo, func(m *tb.Message) {
		b.Bot.Send(b.Member, "Попробуйте еще раз!", menuMember)
		b.Data = nil
		fmt.Println(b.Data)
	})

}

//Цикл
func (b *Bot) endReport() {
	for _, v := range b.Data {
		message := v.Department + "\n" + v.Name + "\n" + v.Date + "\n" + v.Plan
		b.Bot.Send(b.Member, message)
	}
}

//Создание кнопок
func (b *Bot) initMenu() {
	menuMember.Reply(
		menuMember.Row(btnPlan),
	)
	menuData.Reply(
		menuData.Row(btnYes),
		menuData.Row(btnNo),
	)
}

//Авторизация
func (b *Bot) auth() {
	b.Bot.Handle("/start", func(m *tb.Message) {
		b.Member = m.Sender
		b.Bot.Send(b.Member, "Привет сотрудник!", menuMember)
	})
}
func (b *Bot) Init() {
	b.auth()
	b.collector()
	b.initMenu()
	b.checkData()

	b.Bot.Start()
}
