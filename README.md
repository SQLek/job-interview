# wp-interview

Prosty serwer REST z funkcją obserwowania zewnętrznych serwisów http.
Projekt powstał jako zadanie rekrutacyjne.

## Uruchomienie i testy

    docker compose build
    docker compose up
    go test ./...
    go vet ./...
    golangci-lint run ./...

### Uwaga na windowsie

Testy zawarte w tym projekcie wymagają dockera na maszynie.
DockerTest obsługuje windowsa, ale jeszcze nie obsługuje TLS

Prawoklik ikona dockera, ustawienia i załączyć
> Expose daemon on tcp://localhost:2375 without TLS

Takie coś powinno otworzyć nie szyfrowany socket tylko dla lokalnej maszyny.
Nie jestem specem od podsystemu sieciowego windows,
wiec ręki nie dam sobie uciąć.

## Rest API

Pełen opis i specyfikacja w pliku ZadanieGo.pdf.

### GET /api/urls

Listuje adresy które aktualnie są obserwowane.
Zwraca tablicę JSon.

### POST /api/urls

Dodaje adres to cyklicznego obserwowania.
Akceptuje obiekt JSon z polami:

* url - adres który chcemy obserwować
* interval - odstęp czasu pomiedzy odczytami
  nie mniej niż 5 sekund. Domyślnie 60 sekund.

Zwraca oiekt JSon z polem `id` identyfikującym zadania.

### PATCH /api/urls/{ID}

Modyfikuje dane obserwowanego adresu wskazywany przez `ID`.
Akceptuje obiekt JSon z polami:

* url - adres który chcemy obserwować
* interval - odstęp czasu pomiedzy odczytami
  nie mniej niż 5 sekund. Domyślnie 60 sekund.

### DELETE /api/urls/{ID}

Przestaje obserwować dany adres i usuwa wzmianki z bazy danych.

### GET /api/urls/{ID}/history

Listuje historię podglądu danego adresu, jako tablice JSon

## Odpowiedzi na pytania otwarte

### Jakie problemy mogą się pojawić w trakcie działania serwisu?

### W jaki sposób zabezpieczyć serwis przed złośliwymi zapytaniami?

Ten projekt, nawet perfekcyjnie wykonany ma potencjał na nadużycia,
czy wykorzystaniu przez osoby trzecie do DoS.

Z pomysłów jakie mam na zminimalizowanie powierzchni ataku to:

* Dodanie ograniczenia dla kolumny url "unique" aby jeden adres nie mógł się pojawić wielokrotnie
* Normalizacja obserwowanych adresów aby uniknać nadużyć czy pomyłek w rodzaju `example.com?foo=bar` i `example.com?foo=baz`.
* Whitelista z adresami/domenami które przewidujemy obserwować.

### Ile zapytań i danych serwis jest w stanie obsłużyć?

Ile jednoczesnych obserwacji można by było prowadzić na raz?
Więcej niż jedno. Ile dokładnie zależy od maszyny.

Co do zapytań o historię obserwacji dochodzi również jak długą historię system ma przechowywać.

Z pierwszej generacji ryzena, setki zapytań na sekundę możliwe że się wyciśnie.
Przechowywanie obserwacji w bazie `key-value` lub `wide-column` może dociagnie do paru tysiecy na sekundę.

Nie benchmarkowałem tego projektu wiec szacunki proszę uważać że są z kopystką soli.

### W jaki sposób można się przygotować na większy ruch?

Mój pomysł na zwiększenie przepustowości systemu opisałem w `Moja opinia co do serwer/worker`.

### W jaki sposób sprawdzić, że wszystko działa? Co monitorować?

Doświadczenie mam jedynie z prometeuszem/grafaną.
Prometeusz korzysta z tego samego muxa/servera,
co daje nam możliwość sprawdzenia czy kontener żyje.

Monitorujemy ile zapytań http z odpowiedzią 200/400/500 było.
Jeżeli odpowiedzi 400 jest wiecej jak x% można podnieść alarm.
Jeżeli odpowiedzi 500 jest wiecej jak 0.0x% można podnieść alarm.

Można też monitorować QoS (quality of service), na przykład histogramami.
Można obserwować jaki czas reakcji miało 90% zapytań.
Jak wzrośnie powyżej iluś milisekund, podnieść alarm.

## Moje opinie

Luźne opinie odnośnie tego projektu.

### Moje opinie co do API

Trochę dziwi mnie że nie ma paginacji.
Szczególnie w history by się przydała.

Endpoint `UPDATE` powinien według mnie jedynie modyfikować `interval`.
Po modyfikowaniu `url` według mnie dane tracą na użyteczności.

Zastanawia mnie też używanie floatów w jako time.Time i time.Duration.
Realny projekt nie działa w próżni, więc jak jest specyfikacja to się ją trzyma.

## Moja opinia co do serwer/worker

Będąc w zgodzie z specyfikacją można odpalić dwie kopie.
Jedna instancja obsługuje `PUT/PATCH/DELETE` na zadaniach.
Druga instancja obsługuje historię zapytń pod dany URL i listowanie aktualnych zadań.
Zgodne z specyfikacją? Tak. Kiss :*

Przypuszczam że intencją było `pokaż czy potrafisz zrobić mikrousługi`.
Jeżeli tak stawiamy sprawę to:

* Server rozkręca workerów na bieżąco używając docker.sock
  lub innej adekwatnej metody

* Pojedyńczy worker obsługuje jeden adres.
  Takim czymś można byłoby dać dystrybucję obciążenia po róznych serverach.
  Można by było dodać reduntancję.
  SIGTERM przy usuwaniu adresu.
  SIGHUP przy modyfikacji zadania.

* Zbiór wyników umieściłbym w osobnej bazie danych.
  Parę dni temu robiłem zadanie rekrutacyjne z Cassandrą i była by to fajna opcja.
  Można byłoby rozważyć Redis, Kafka czy Influx ale nie jestem biegły w tych bazach.

* Ruch dodaj/usuń zadanie zostawiłbym na bazie SQL,
  ponieważ nie sądzę aby te czynności były wykonywane często.
  Nie wiem jak MySQL ale postgresa można shardować w stylu master-slaves dla reduntancji.

To jest niestety więcej jak parę dni którymi dysponowałem.

## Dług techniczny

W moich projektach utrzymuję listę
`nie do końca eleganckie, działa, do poprawki jak będzie luźniej`
Oto lista co można by było proprawić, czy przemyśleć w wolnej chwili.

### GORM

To jest mój pierwszy projekt z bibloteką GORM.
Nie jestem biegły z MySQL wiec zdecydowałem się ugryźć temat ORM-em.
Bibloteka sprawuje sięcałkiem fajnie,
co nie zmienia faktu że większe rozeznanie z nią, by mi nie zaszkodziło.

Jak embeduje się gorm.Model to delete nie jest twarde i usuniete dane by zostawały w bazie.
Bez osadzenia gorm.Model, wbudowane createdAt przestało działać.

Dodałbym konteksty i timeouty do zapytań listTasks/listEntries.
Może mimo wszystko paginacja.
Do prototypu nie widzę potrzeby.

### DockerTest

DockerTest bardzo przydatne narzedzie.
Zamiast mockować można testować z prawdziwą bazą i prawdziwym driverem.

Dotychczas radziłem sobie przy pomocy:
    docker compose build
    docker compose up
    go test ./...
    docker compose down

Z DockerTest, w kodzie szykujacym testy,
można odpalać kontenery z usługami potrzebnymi przez testy.

Baza danych zawsze czysta.
Nie trzeba ręcznie przydzielać portów.

Zamierzam poswiecić tematowi więcej czasu w przyszłosci.

## Zakończenie

Dziękuję za poświecony czas w audyt mojego rozwiązania na zadanie rekrutacyjne.
