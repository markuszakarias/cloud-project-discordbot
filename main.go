package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
)

// api key for news letter: cfa7f832f70e41c899bf6b735ef77abf

var exampleResponse = `{
    "status": "ok",
    "totalResults": 30,
    "articles": [
        {
            "source": {
                "id": null,
                "name": "Digi.no"
            },
            "author": "Odd Richard Valmot",
            "title": "Samsung blir en seriøs utfordrer i PC-markedet - digi.no",
            "description": "Lettere, tynnere, store batterier, AMOLED skjermer, lavere priser og kobling til mobilene skal posisjonere Samsung i PC-markedet",
            "url": "https://www.digi.no/artikler/samsung-blir-en-serios-utfordrer-i-pc-markedet/509686",
            "urlToImage": "https://img.gfx.no/2672/2672321/Galaxy_Book_Pro_360_13inch_MysticBronze_S_Pen_Broll_2_210417012109.1200x676.jpg",
            "publishedAt": "2021-04-28T14:00:00Z",
            "content": "For et års tid siden lanserte Samsung PCer i Norge. Flotte PCer, men dyre og kanskje litt halvhjertet markedsført. Pussig nok hadde ingen av dem trådløst bredbånd, noe vi burde forvente av en stor mo… [+4714 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Dagbladet.no"
            },
            "author": "Emma Cecilia Eriksson",
            "title": "Vladislav Ivanov: - «Fanget» i realityshow - endelig stemt ut - Dagbladet.no",
            "description": "Var «fanget» i realityshow i tre måneder, etter at fans nektet å stemme ham ut.",
            "url": "https://www.dagbladet.no/kjendis/fanget-i-realityshow---endelig-stemt-ut/73697401",
            "urlToImage": "https://www.dagbladet.no/images/73699521.jpg?imageId=73699521&panow=100&panoh=100&panox=0&panoy=0&heightw=100&heighth=100&heightx=0&heighty=0&width=1200&height=630",
            "publishedAt": "2021-04-28T13:45:01Z",
            "content": "Beslutningen om å melde seg på et realityprogram bør kanskje ikke tas hvis du ikke virkelig er skråsikker. Det har nå en russisk mann fått erfare, etter at han ble «fanget» i en boyband-konkurranse i… [+4416 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Abcnyheter.no"
            },
            "author": "David Stenerud",
            "title": "Helt riktig å slippe til Charter-Svein i Debatten på NRK - ABC Nyheter",
            "description": "NRK Debatten ga Svein Østvik anledning til å tale sin sak. Heldigvis gikk det rett vest.",
            "url": "https://www.abcnyheter.no/stemmer/2021/04/28/195755405/helt-riktig-a-slippe-til-charter-svein-i-debatten-pa-nrk",
            "urlToImage": "https://imaginary.abcmedia.no/resize?width=980&interlace=true&url=https%3A%2F%2Fpresizely.abcmedia.no%2Fhttps%3A%2F%2Fabcnyheter.drpublish.aptoma.no%2Fout%2Fimages%2Farticle%2F%2F2021%2F04%2F28%2F195755423%2F1%2Foriginal%2F43866261.jpg",
            "publishedAt": "2021-04-28T13:44:31Z",
            "content": "Hvis motstanderen din er i ferd med å grave sin egen grav, ikke slåss med ham om spaden, er det visst noe som heter. Det er litt sånn jeg har det med «Charter-Svein» Østvik og hans meningsfeller. En … [+2523 chars]"
        },
        {
            "source": {
                "id": "aftenposten",
                "name": "Aftenposten"
            },
            "author": "Odd Inge Aas, Alf Ole Ask",
            "title": "Velgerne mener KrF kun er best på ett område. Er det nok til å berge partiet fra undergangen? - Aftenposten",
            "description": "Fire og en halv måneder før valget er KrF godt under sperregrensen.",
            "url": "https://www.aftenposten.no/norge/politikk/i/vAJBkl/velgerne-mener-krf-kun-er-best-paa-ett-omraade-er-det-nok-til-aa-berge-p",
            "urlToImage": "https://premium.vgc.no/v2/images/76aa1c3d-6373-45ee-a746-b47c79021637?fit=crop&h=1072&w=2048&s=d95dfb2423ac593cb3e41467c7918e31b258f641",
            "publishedAt": "2021-04-28T13:30:25Z",
            "content": "Fire og en halv måneder før valget er KrF godt under sperregrensen.\r\nKrF-leder Kjell Ingolf Ropstad Foto: Signe Dons\r\nSist oppdatert nå nettopp\r\nKjell Ingolf Ropstad er klar over alvoret. I helgen sk… [+335 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Www.vg.no"
            },
            "author": "Ole Kristian Strøm, Joachim Baardsen",
            "title": "Godset-styret anmelder saken om Henrik Pedersen til NFF - VG",
            "description": "Styret i Strømsgodset mener at det er grunn til å konkludere med at eks-trener Henrik Pedersen brukte uakseptable ord og uttrykk. De anmelder saken til NFF.",
            "url": "https://www.vg.no/sport/fotball/i/7KzrK9/godset-styret-anmelder-saken-om-henrik-pedersen-til-nff",
            "urlToImage": "https://akamai.vgc.no/v2/images/3bb225ee-5950-4f62-b443-a5c3743960f7?fit=crop&h=1140&w=1900&s=b20ae0cf77a7f8199f15b8d0d08812c979261950",
            "publishedAt": "2021-04-28T13:28:32Z",
            "content": "EKS-TRENER: Henrik Pedersen fotografert etter en kamp mot Rosenborg i november 2019. Foto: Berit Roald\r\nStyret i Strømsgodset mener at det er grunn til å konkludere med at eks-trener Henrik Pedersen … [+5737 chars]"
        },
        {
            "source": {
                "id": "nrk",
                "name": "NRK"
            },
            "author": "Fredrik Moen Gabrielsen, Ugo Fermariello, Svein Vestrum Olsson, Simen Hunding Strømme, Joakim Reigstad",
            "title": "Fullvaksinerte kan ha mer kontakt privat - NRK",
            "description": "Fullvaksinerte kan ha nær kontakt med andre fullvaksinerte og uvaksinerte utenfor risikogruppen. Men dette gjelder bare i private hjem. I det offentlige rom blir det ingen endringer.",
            "url": "https://www.nrk.no/norge/fullvaksinerte-kan-ha-mer-kontakt-privat-1.15473502",
            "urlToImage": "https://gfx.nrk.no/3itHHEZ-5jmBhh-OC4LYsQqxgLS_eN2tQW-yw4MibQ7A.jpg",
            "publishedAt": "2021-04-28T13:18:23Z",
            "content": "Denne artikkelen er over en måned gammel, og kan inneholde utdaterte råd fra myndighetene angående koronasmitten.\r\nHold deg oppdatert i NRKs oversikt, eller gjennom FHIs nettsider.\r\nHøie sier de har … [+8228 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Www.tek.no"
            },
            "author": null,
            "title": "Komplett med ny PS5-dato: 31. desember - Tek.no - Tek.no",
            "description": "Norges grundigste tester, guider og nyheter relatert til forbrukerteknologi finner du på Tek.no. Med langt over en million brukere i uken er ingen større enn oss.",
            "url": "https://www.tek.no/nyheter/nyhet/i/VqJAEp/komplett-med-ny-ps5-dato-31-desember",
            "urlToImage": "https://cdn.tek.no/pro/no/current/assets/logos/fallback-og-image-2.png",
            "publishedAt": "2021-04-28T13:04:31Z",
            "content": null
        },
        {
            "source": {
                "id": null,
                "name": "Www.tv2.no"
            },
            "author": "TV 2 AS",
            "title": "Solskjær svarer Roma-fansen: - Ikke meningen å være respektløs - TV 2",
            "description": "Ole Gunnar Solskjær håper å bryte semifinale-forbannelsen og svarer Roma-fansen.",
            "url": "https://www.tv2.no/a/13970807/",
            "urlToImage": "https://www.cdn.tv2.no/images/13968676.jpg?imageId=13968676&panow=100&panoh=100&panox=0&panoy=0&heightw=100&heighth=100&heightx=0&heighty=0&width=1200&height=630",
            "publishedAt": "2021-04-28T12:56:37Z",
            "content": "Fire semifinaler, fire nederlag. Det er neppe en statistikk Ole Gunnar Solskjær er særlig fornøyd med, men det er altså fasit etter kristiansunderens første fire semifinaler som Manchester United-sje… [+3364 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Nettavisen.no"
            },
            "author": "Martin Busk",
            "title": "Roma-fansen forbannet på Solskjær - slik svarer United-manageren - Nettavisen",
            "description": "Ole Gunnar Solskjær sier han ikke mente å fornærme noen før semifinalen i Europa League.",
            "url": "https://www.nettavisen.no/sport/roma-fansen-forbannet-pa-solskjar-slik-svarer-united-manageren/s/12-95-3424120597",
            "urlToImage": "https://g.api.no/obscura/API/dynamic/r1/nadp/tr_2000_2000_s_f/1619614890321/2021/04/28/3424120617/1/original/43870603.jpg?chk=7C112F",
            "publishedAt": "2021-04-28T12:45:10Z",
            "content": "Ole Gunnar Solskjær sier han ikke mente å fornærme noen før semifinalen i Europa League.\r\n28.04.21 14:45\r\n28.04.21 15:01\r\nDet er Ole Gunnar Solskjærs uttalelser om Roma for et par dager siden som har… [+2263 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Www.vg.no"
            },
            "author": "Simon Zetlitz Nessler",
            "title": "Ruud enkelt videre i München: – Spilte en god kamp - VG",
            "description": "(Casper Ruud – Pablo Cuevas 6–3, 6–2) Casper Ruud (22) brukte bare en drøy time på å ta seg til kvartfinalen i BMW Open i München. Nå er mulighetene svært gode for å ta seg til en finale i grusturneringen.",
            "url": "https://www.vg.no/sport/i/oAOroW/ruud-enkelt-videre-i-munchen-spilte-en-god-kamp",
            "urlToImage": "https://akamai.vgc.no/v2/images/e21442bc-a579-4763-82d5-d056cecfbd01?fit=crop&h=1215&w=1900&s=d70e5903e920318ca3663b8f0018e5218af9ccb9",
            "publishedAt": "2021-04-28T12:20:41Z",
            "content": "KLAR FOR KVARTFINALE: Casper Ruud vant i strake sett mot Pablo Cuevas. Foto: CHRISTOF STACHE / AFP\r\n(Casper Ruud Pablo Cuevas 63, 62) Casper Ruud (22) brukte bare en drøy time på å ta seg til kvartfi… [+3220 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Www.vg.no"
            },
            "author": "Christina Quist, Hanna Haug Røset, Vilde Elgaaen, Gordon Andersen, Sven Arne Buggeland, Ådne Husby Sandnes, Morten S. Hopperstad, Terje Bringedal (foto), Fredrik Solstad (foto), Tore Kristiansen (foto), Helge Mikalsen (foto)",
            "title": "Drapssiktet skylder drept kvinne over 12 millioner etter byggekonflikt - VG",
            "description": "Den drapssiktede mannen i slutten av 30-årene skal ha vært i en årelang økonomisk strid med kvinnen politiet nå mener han har skutt og drept.",
            "url": "https://www.vg.no/nyheter/innenriks/i/41kpWo/drapssiktet-skylder-drept-kvinne-over-12-millioner-etter-byggekonflikt",
            "urlToImage": "https://akamai.vgc.no/v2/images/061d992f-5422-4d22-b562-f3c1adfebc06?fit=crop&h=1267&w=1900&s=de35d5e5fa849e7e7e1ffa46d45c6d623b7866b4",
            "publishedAt": "2021-04-28T12:15:29Z",
            "content": "STORE STYRKER: Store politistyrker rykket ut til stedet etter at det kom melding om at en person var skutt og drept på åpen gate på Frogner. Foto: Fredrik Solstad, VG\r\nDen drapssiktede mannen i slutt… [+4140 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Seher.no"
            },
            "author": null,
            "title": "Justin Biebers slaktes for sin nye hårfrisyre - Seoghør.no",
            "description": "Superstjernens nye hårfrisyre har fått hard medfart.",
            "url": "https://www.seher.no/kjendis/slaktes-for-nytt-har/73697746",
            "urlToImage": "https://www.seher.no/images/73697395.jpg?imageId=73697395&panow=100&panoh=33.971291866029&panox=0&panoy=14.832535885167&heightw=60.810810810811&heighth=100&heightx=28.378378378378&heighty=0&width=1200&height=630",
            "publishedAt": "2021-04-28T12:09:56Z",
            "content": "Justin Bieber (27). Unggutten, med den silkemyke stemmen, som en gang tok verden med storm, men som de seneste årene har stått for mang en skandale.\r\nMusikkvideoen for låten «Baby» er en av de mest m… [+3247 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Www.vg.no"
            },
            "author": "Catherine Gonsholt Ighanian",
            "title": "Rachel Bilson og Rami Malek har skværet opp - VG",
            "description": "Rachel Bilson (39) fortalte hele verden at Rami Malek (39) ble sur da hun postet et gammelt bilde av de to.",
            "url": "https://www.vg.no/rampelys/film/i/qA49jo/rachel-bilson-og-rami-malek-har-skvaeret-opp",
            "urlToImage": "https://akamai.vgc.no/v2/images/c9f20fa3-0a17-4156-8fae-35405187829d?fit=crop&h=1200&w=1449&s=044fcd4e0e565713c8046a89999b8bc2795e07dd",
            "publishedAt": "2021-04-28T12:06:43Z",
            "content": "STUDIEKAMERATER: Rami Malek og Rachel Bilson, her avbildet på hver sin kjendisfest i henholdsvis 2019 og 2018. Foto: Pa Photos og AFP\r\nRachel Bilson (39) fortalte hele verden at Rami Malek (39) ble s… [+1806 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "E24.no"
            },
            "author": "Ine Brunborg, Joachim Birger Nilsen, Infront TDN Direkt",
            "title": "Oljeprisen stiger - oppgang på Oslo Børs - E24",
            "description": "Børsen steg onsdag etter morgenens resultatslipp. Storebrand falt, mens Aker BP endte opp. Oljeprisen stiger etter nedgang tidligere på dagen.",
            "url": "https://e24.no/boers-og-finans/i/Qm4vnJ/oppgang-paa-oslo-boers",
            "urlToImage": "https://smp.vgc.no/v2/images/fff8e81b-038a-4cc8-b68b-bb69ba793a78?fit=crop&h=1267&w=1900&s=b88dbb3e008799ab530129da72d6c14b5500d72d",
            "publishedAt": "2021-04-28T12:03:51Z",
            "content": "Børsen steg onsdag etter morgenens resultatslipp. Storebrand falt, mens Aker BP endte opp. Oljeprisen stiger etter nedgang tidligere på dagen.\r\nOslo Børs er inne i resultatsesongen for selskapenes fø… [+2992 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Www.bt.no"
            },
            "author": "Roar Lyngøy",
            "title": "Smitteutbrudd på Haakonsvern - alle permisjoner inndras abonnent - Bergens Tidende",
            "description": "– Det er en utrolig trasig situasjon, sier Michel Hayes, talsperson for Sjøforsvaret.",
            "url": "https://www.bt.no/nyheter/lokalt/i/BlPQo7/smitteutbrudd-paa-haakonsvern-alle-permisjoner-inndras",
            "urlToImage": "https://premium.vgc.no/v2/images/efc5c01f-1598-4468-bcc3-16cb6e7adec9?fit=crop&h=995&w=1900&s=41a0a88dc3dd91907bed74bbe00035d761ab2a0a",
            "publishedAt": "2021-04-28T11:48:31Z",
            "content": "Det er en utrolig trasig situasjon, sier Michel Hayes, talsperson for Sjøforsvaret.\r\nPublisert Publisert For mindre enn 30 minutter siden\r\nLes hele saken med abonnement\r\nAllerede abonnent? Logg inn"
        },
        {
            "source": {
                "id": "nrk",
                "name": "NRK"
            },
            "author": "Øyvind Nyborg",
            "title": "Boris skal granskes - NRK",
            "description": "Krisen fortsetter for Boris Johnson. Valgkommisjonen har satt i gang full etterforskning av hvem som betalte for oppussingen av leiligheten hans i Downing Street.",
            "url": "https://www.nrk.no/urix/boris-skal-granskes-1.15474145",
            "urlToImage": "https://gfx.nrk.no/HuOLtTO0u2RXd1U7P11PvAtO3QibhYVhxNnIrjxhuxQQ.jpg",
            "publishedAt": "2021-04-28T11:48:01Z",
            "content": "I en uttalelse fra Valgkommisjonen heter det at det er grunn til å tro at det har skjedd ulovligheter og at det er derfor Johnson nå skal ettergås. Dersom Johnson blir funnet skyldig i å ha tatt imot… [+1933 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "E24.no"
            },
            "author": "Fabian Skalleberg Nilsen",
            "title": "Spotify snur fra minus til pluss - E24",
            "description": "Strømmekjempen startet 2021 med en opptur, i både kroner og ører. Men Spotify-aksjen raser på Wall Street etter nedjusterte anslag om brukerveksten.",
            "url": "https://e24.no/internasjonal-oekonomi/i/JJemg7/spotify-stuper-paa-boers-til-tross-for-resultatbedring",
            "urlToImage": "https://smp.vgc.no/v2/images/528f8d53-25c6-4acf-a0d2-4a5fb0264ac0?fit=crop&h=1267&w=1900&s=341ee68ffbd7088093710d928fa04b97a70f674a",
            "publishedAt": "2021-04-28T11:36:53Z",
            "content": "Strømmekjempen startet 2021 med en opptur, i både kroner og ører. Men Spotify-aksjen raser på Wall Street etter nedjusterte anslag om brukerveksten.\r\nKLANGBUNN: 2020 ble et tøft år for mange, mens år… [+2362 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Seher.no"
            },
            "author": null,
            "title": "Öde Nerdrum: - Slik svarer han på kjærestespørsmål - Seoghør.no",
            "description": "Det er uklart hva Öde Nerdrums sivilstatus er i dag.",
            "url": "https://www.seher.no/kjendis/slik-svarer-han-pa-kjaerestesporsmal/73698749",
            "urlToImage": "https://www.seher.no/images/73698292.jpg?imageId=73698292&panow=100&panoh=50.714285714286&panox=0&panoy=0&heightw=100&heighth=100&heightx=0&heighty=0&width=1200&height=630",
            "publishedAt": "2021-04-28T10:52:24Z",
            "content": "Kunstneren Öde Nerdrum (26) tok seg helt til finaleuka i vinterens utgave av «Farmen Kjendis», og ble en stor seerfavoritt. Han vant ikke konkurransen, men med sitt rufsete hår, riksmålstale og sjarm… [+3370 chars]"
        },
        {
            "source": {
                "id": "nrk",
                "name": "NRK"
            },
            "author": "Fredrik Moen Gabrielsen, Svein Vestrum Olsson",
            "title": "Vaksinerte Tone: – Håper det blir litt mer frihet - NRK",
            "description": "Onsdag kommer det nye råd for dem som er vaksinert i Norge. Vaksinerte Tone Walløe Strøm (67) håper på mer frihet og lengter etter reising og å treffe flere av vennene sine.",
            "url": "https://www.nrk.no/norge/vaksinerte-tone_-_-haper-det-blir-litt-mer-frihet-1.15473942",
            "urlToImage": "https://gfx.nrk.no/Ta-C2zMAPaue0kizv_VoKAFMKBdejEWpEB-oPDV1IImA.jpg",
            "publishedAt": "2021-04-28T09:53:51Z",
            "content": "Denne artikkelen er over en måned gammel, og kan inneholde utdaterte råd fra myndighetene angående koronasmitten.\r\nHold deg oppdatert i NRKs oversikt, eller gjennom FHIs nettsider.\r\nJeg håper det bli… [+4145 chars]"
        },
        {
            "source": {
                "id": null,
                "name": "Www.vg.no"
            },
            "author": "Frank Ertesvåg",
            "title": "Ap borer i manglende vaksinering i skoler - Melby venter på FHI - VG",
            "description": "Aps utdanningspolitiker Torstein Tvedt Solberg grillet kunnskapsministeren om hvorfor vaksineringen har uteblitt for skole- og barnehageansatte. Guri Melby svarer at hun avventer råd fra FHI i midten av mai.",
            "url": "https://www.vg.no/nyheter/innenriks/i/mBlo74/ap-borer-i-manglende-vaksinering-i-skoler-melby-venter-paa-fhi",
            "urlToImage": "https://akamai.vgc.no/v2/images/94d73ba6-c731-498f-b6d1-f74d5af95e38?fit=crop&h=1190&w=1900&s=d2d953020a0348dbbd2fb25c96a99b240b31a3e3",
            "publishedAt": "2021-04-28T09:42:48Z",
            "content": "LÆRER-STIKK: Arbeiderpartiet vil prioritere vaksinering av personell i skoler og barnehager. ILLUSTRASJONSFOTO Foto: Berit Roald, NTB\r\nAps utdanningspolitiker Torstein Tvedt Solberg grillet kunnskaps… [+3805 chars]"
        }
    ]
}`

// example request
//
// https://newsapi.org/v2/top-headlines?country=no&apiKey=API_KEY
//

type newsLetter struct {
	author         string `json:"author"`
	date_published string `json:"date_published"`
	title          string `json:"title"`
	description    string `json:"description"`
	url_to_story   string `json:"url_to_story"`
}

type newsLetters struct {
	newsletters []newsLetter
}

func populateNewsLetters(paramStruct newsLetter, jsonResponseString string) newsLetter {
	test1 := gjson.Get(jsonResponseString, "articles.1.author")
	test2 := gjson.Get(jsonResponseString, "articles.1.publishedAt")
	test3 := gjson.Get(jsonResponseString, "articles.1.title")
	test4 := gjson.Get(jsonResponseString, "articles.1.description")
	test5 := gjson.Get(jsonResponseString, "articles.1.url")

	test1String := test1.String()
	test2String := test2.String()
	test3String := test3.String()
	test4String := test4.String()
	test5String := test5.String()

	paramStruct.author = test1String
	paramStruct.date_published = test2String
	paramStruct.title = test3String
	paramStruct.description = test4String
	paramStruct.url_to_story = test5String
	return paramStruct
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	url := "https://newsapi.org/v2/top-headlines?country=no&apiKey=cfa7f832f70e41c899bf6b735ef77abf"

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
	}

	r.Header.Add("content-type", "application/json")

	client := &http.Client{}

	// Issue request
	res, err := client.Do(request)
	if err != nil {
		fmt.Errorf("Error in response:", err.Error())
	}

	// Print output
	output, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Errorf("Error when reading response: ", err.Error())
	}

	jsonResponseAsString := string(output)

	var test newsLetter
	articleStruct := populateNewsLetters(test, jsonResponseAsString)
	fmt.Printf("%+v\n", articleStruct)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", simpleHandler)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
