package app

import (
	"WB_Tech_level_0/internal/cache"
	"WB_Tech_level_0/internal/models"
	"WB_Tech_level_0/internal/store"
	natss "WB_Tech_level_0/internal/transport/nats"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
)

type Service struct {
	config     *Config
	natsserver *natss.NatsServer
	cache      *cache.Cache
	store      *store.Store
	router     *mux.Router
}

// Инициализация сервиса
func New(config *Config) *Service {
	return &Service{
		config: config,
		router: mux.NewRouter(),
	}
}

// Запуск сервиса
func (s *Service) Start() error {
	if err := s.configureNats(); err != nil {
		return err
	}
	if err := s.configureDataBase(); err != nil {
		return err
	}
	defer s.store.Close()
	s.configureCache()
	if err := s.recdb(); err != nil {
		fmt.Println("Восстановление данные из бд не удалось:", err)
	}
	s.StartNats()
	defer s.natsserver.Close()
	s.configureRouter()

	http.ListenAndServe(s.config.BindAddr, s.router)
	return nil
}

// Определение точек входа для роутинга
func (s *Service) configureRouter() {
	s.router.HandleFunc("/", s.indexId())
	s.router.HandleFunc("/postform", s.getId())
}

// Обработчик для Роутера - выдача товара по id
func (s *Service) getId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, exist := s.cache.GetCache(r.FormValue("id"))
		if !exist {
			log.Fatal("element doesn't exist")
		}
		fmt.Println(data.Items)
		tmpl, err := template.ParseFiles("web/getid.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, data)
	}
}

// Обработчик для Роутера - точка входа
func (s *Service) indexId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("web/index.html")
		if err != nil {
			log.Fatal(err)
		}
		tmpl.Execute(w, nil)
	}
}

// Восстановление данных из БД в КЭШ
func (s *Service) recdb() error {
	data, err := s.store.Message().Read_ALL()
	if err != nil {
		return err
	}
	for _, v := range data {
		if err := s.cache.InCache(v); err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("Восстановление из БД прошло успешно")
	return nil
}

// Инициализация КЭШа
func (s *Service) configureCache() {
	s.cache = cache.New()
}

// Передача конфига и создание БД
func (s *Service) configureDataBase() error {
	db := store.New(s.config.Store)
	if err := db.Open(); err != nil {
		return err
	}
	s.store = db

	return nil
}

// Передача конфига и создание NATS
func (s *Service) configureNats() error {
	ns := natss.New(s.config.Nats)
	if err := ns.InitConnect(); err != nil {
		return err
	}
	s.natsserver = ns
	return nil
}

// Обработчик NATS
func (s *Service) StartNats() {
	s.natsserver.Natc.Subscribe(s.config.Nats.Subject, func(msg *nats.Msg) {
		var nats_data models.DataJson
		if err := json.Unmarshal(msg.Data, &nats_data); err != nil {
			log.Fatal(err)
		}
		data := convertdata(&nats_data)
		if err := s.cache.InCache(data); err != nil {
			fmt.Println(err)
		} else {
			err := s.store.Message().Create(&data)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Произведена запись данных в БД")
			}
		}
	})
}
func convertdata(d *models.DataJson) models.Data {
	data := models.Data{}
	data.Order.Customer_id = d.Customer_id
	data.Order.Delivery_service = d.Delivery_service
	data.Order.Entry = d.Entry
	data.Order.Date_created = d.Date_created
	data.Order.Internal_signature = d.Internal_signature
	data.Order.Locale = d.Locale
	data.Order.Oof_shard = d.Oof_shard
	data.Order.Shardkey = d.Shardkey
	data.Order.Track_number = d.Track_number
	data.Order.Sm_id = d.Sm_id
	data.Order.Order_uid = d.Order_uid
	data.Delivery = d.Delivery
	data.Items = d.Items
	data.Payment = d.Payment
	return data
}
