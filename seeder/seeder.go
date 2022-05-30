package main

import (
	"context"
	"database/sql"
	"log"

	sqlc "gubuk-service/db/sqlc"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var testDB *sql.DB
var testQueries *sqlc.Queries

var ownerID = uuid.New()

func main() {
	db, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/gubukid?sslmode=disable")
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	testDB = db
	testQueries = sqlc.New(db)

	seedUser()
	seedHome()
}

func seedUser() {
	var err error
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	// user without avatar, with role as a tenant
	_, err = testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		ID:          uuid.New(),
		Fullname:    "Febrian Amir",
		Username:    "febrian",
		Email:       "febrian@gmail.com",
		Role:        "tenant",
		Gender:      "male",
		PhoneNumber: "0812345678",
		Password:    string(hashedPassword),
		Address:     "Jln Bollangi",
		Avatar:      "",
	})
	if err != nil {
		log.Fatal(err)
	}

	// user with avatar, with role as an owner
	_, err = testQueries.CreateUser(context.TODO(), sqlc.CreateUserParams{
		ID:          ownerID,
		Fullname:    "Amiruddin",
		Username:    "amiruddin",
		Email:       "amiruddin@gmail.com",
		Role:        "owner",
		Gender:      "male",
		PhoneNumber: "0812345678",
		Password:    string(hashedPassword),
		Address:     "Jln Bollangi",
		Avatar:      "https://res.cloudinary.com/quikzens/image/upload/v1652707630/avatar/dfvulizzdcmqydlynet3.png",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func seedHome() {
	var err error
	_, err = testQueries.CreateHouse(context.TODO(), sqlc.CreateHouseParams{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Title:         "Rumah Krong Bade",
		FeaturedImage: "https://res.cloudinary.com/quikzens/image/upload/v1653120218/house/krong-bade_ixgjms.jpg",
		Bedrooms:      1,
		Bathrooms:     1,
		TypeRent:      "day",
		Price:         100000,
		ProvinceID:    11,
		CityID:        1101,
		Description:   "Rumah Krong Bade dari Aceh ini berbentuk memanjang dari timur ke barat menyerupai persegi panjang. Di bagian depan rumah dilengkapi dengan tangga untuk masuk ke dalam rumah.",
		Amenities:     "Furnished",
		Area:          70,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = testQueries.CreateHouse(context.TODO(), sqlc.CreateHouseParams{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Title:         "Rumah Bolon",
		FeaturedImage: "https://res.cloudinary.com/quikzens/image/upload/v1653120218/house/bolon_veosjm.jpg",
		Bedrooms:      2,
		Bathrooms:     2,
		TypeRent:      "month",
		Price:         4000000,
		ProvinceID:    12,
		CityID:        1201,
		Description:   "Pada rumah adat Bolon ini, terdapat dua bagian yang berbeda, yaitu Jabu Bolon dan juga Jabu Parsakitan. Jabu Bolon biasa menjadi tempat untuk keluarga besar, sedangkan Jabu Parsakitan adalah tempat untuk membicarakan masalah adat.",
		Amenities:     "Shared Accomodation",
		Area:          60,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = testQueries.CreateHouse(context.TODO(), sqlc.CreateHouseParams{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Title:         "Rumah Gadang",
		FeaturedImage: "https://res.cloudinary.com/quikzens/image/upload/v1653120217/house/gadang_ijq7m4.jpg",
		Bedrooms:      3,
		Bathrooms:     3,
		TypeRent:      "year",
		Price:         70000000,
		ProvinceID:    13,
		CityID:        1301,
		Description:   "Rumah adat Gadang terbuat dari ijuk dan bentuknya mirip seperti tanduk kerbau, yang melambangkan kemenangan suku Minang dalam perlombaan adu kerbau di Jawa.",
		Amenities:     "Pet Allowed",
		Area:          54,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = testQueries.CreateHouse(context.TODO(), sqlc.CreateHouseParams{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Title:         "Rumah Selaso Jatuh Kembar",
		FeaturedImage: "https://res.cloudinary.com/quikzens/image/upload/v1653120219/house/selasa-jatuh-kembar_omk3uo.jpg",
		Bedrooms:      4,
		Bathrooms:     4,
		TypeRent:      "day",
		Price:         50000,
		ProvinceID:    14,
		CityID:        1401,
		Description:   "Rumah ini memiliki arti rumah dengan dua selasar. Masyarakat Riau tidak menjadikan Rumah Selaso Jatuh Kembar sebagai tempat tinggal mereka, tetapi hanya menggunakannya untuk acara adat.",
		Amenities:     "Furnished,Pet Allowed",
		Area:          45,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = testQueries.CreateHouse(context.TODO(), sqlc.CreateHouseParams{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Title:         "Rumah Bubungan Lima",
		FeaturedImage: "https://res.cloudinary.com/quikzens/image/upload/v1653120217/house/bubungan-lima_xgrdn8.jpg",
		Bedrooms:      5,
		Bathrooms:     5,
		TypeRent:      "month",
		Price:         8500000,
		ProvinceID:    17,
		CityID:        1701,
		Description:   "Rumah adat dari Bengkulu ini memiliki tiang penopang dan menggunakan kayu khusus untuk membuatnya, yaitu kayu Medang Kemuning. Untuk memasuki rumah ini, Anda juga harus menggunakan tangga, yang berada pada bagian depan rumah. ",
		Amenities:     "Furnished,Shared Accomodation",
		Area:          70,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = testQueries.CreateHouse(context.TODO(), sqlc.CreateHouseParams{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Title:         "Rumah Panggung",
		FeaturedImage: "https://res.cloudinary.com/quikzens/image/upload/v1653120219/house/panggung_ebnm1z.jpg",
		Bedrooms:      1,
		Bathrooms:     1,
		TypeRent:      "year",
		Price:         75000000,
		ProvinceID:    15,
		CityID:        1501,
		Description:   "Orang-orang sering menyebutkan bagian atap dari Rumah Panggung ini sebagai “Gajah Mabuk” karena bentuknya yang menyerupai perahu dengan ujung melengkung. Biasanya, rumah adat dari Jambi digunakan untuk tempat tinggal dan juga tempat bermusyawarah.",
		Amenities:     "Shared Accomodation,Pet Allowed",
		Area:          60,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = testQueries.CreateHouse(context.TODO(), sqlc.CreateHouseParams{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Title:         "Rumah Nuwo Sesat",
		FeaturedImage: "https://res.cloudinary.com/quikzens/image/upload/v1653120218/house/nuwo-sesat_o5fzb0.jpg",
		Bedrooms:      2,
		Bathrooms:     2,
		TypeRent:      "day",
		Price:         60000,
		ProvinceID:    18,
		CityID:        1801,
		Description:   "Rumah adat Provinsi Lampung memiliki nama Nuwo Sesat. Ciri khas dari rumah ini adalah bentuknya panggung dan di sisi-sisinya terdapat ornamen yang khas. Biasanya, ukuran dari rumah ini sangat besar, tetapi saat ini banyak yang membuat Rumah Nuwo Sesat berukuran lebih kecil.",
		Amenities:     "Furnished,Pet Allowed,Shared Accomodation",
		Area:          54,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = testQueries.CreateHouse(context.TODO(), sqlc.CreateHouseParams{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Title:         "Rumah Limas",
		FeaturedImage: "https://res.cloudinary.com/quikzens/image/upload/v1653120217/house/limas_rhzhx7.jpg",
		Bedrooms:      3,
		Bathrooms:     3,
		TypeRent:      "month",
		Price:         4500000,
		ProvinceID:    16,
		CityID:        1601,
		Description:   "Rumah adat satu ini memiliki bentuk yang sesuai dengan namanya, yaitu menyerupai limas. Tamu yang berkunjung ke rumah ini harus singgah ke ruang atas atau teras rumah. Hal ini merupakan tradisi masyarakat Sumatera Selatan agar dapat merasakan budaya mereka, yang tampak pada ukiran di dalamnya.",
		Amenities:     "Furnished,Pet Allowed,Shared Accomodation",
		Area:          45,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = testQueries.CreateHouse(context.TODO(), sqlc.CreateHouseParams{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Title:         "Rumah Rakit",
		FeaturedImage: "https://res.cloudinary.com/quikzens/image/upload/v1653120219/house/rakit_rgrfou.jpg",
		Bedrooms:      4,
		Bathrooms:     4,
		TypeRent:      "year",
		Price:         59000000,
		ProvinceID:    19,
		CityID:        1901,
		Description:   "Karena Bangka Belitung memiliki banyak yang tergenang air atau di tepi laut, warga setempat harus menyesuaikan diri, yaitu dengan membangun rumah di atas air juga yang dinamakan Rumah Rakit. Bentuk rumah adat provinsi Bangka belitung terlihat sangat unik karena merupakan perpaduan rumah Melayu dengan aksen arsitektur Tionghoa.",
		Amenities:     "Furnished,Pet Allowed,Shared Accomodation",
		Area:          70,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = testQueries.CreateHouse(context.TODO(), sqlc.CreateHouseParams{
		ID:            uuid.New(),
		OwnerID:       ownerID,
		Title:         "Rumah Atap Limas Potong",
		FeaturedImage: "https://res.cloudinary.com/quikzens/image/upload/v1653120217/house/atap-limas-potong_pbktmd.jpg",
		Bedrooms:      5,
		Bathrooms:     5,
		TypeRent:      "day",
		Price:         90000,
		ProvinceID:    14,
		CityID:        1401,
		Description:   "Rumah adat dari Kepulauan Riau ini terlihat sangat sederhana. Berbentuk seperti rumah panggung, yang memanjang ke belakang dengan dinding kayu tersusun secara vertikal.",
		Amenities:     "Furnished,Shared Accomodation",
		Area:          54,
	})
	if err != nil {
		log.Fatal(err)
	}
}
