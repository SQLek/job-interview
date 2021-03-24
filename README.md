# wp-interview

Zadanie rekrutacyjne dla Wirtualnej Polski.

Zadanie jest aktualnie dopinane, weryfikowane i wysyłane.
Jest udostępnione w takiej formie aby zrobić uprawnienia na githubie.

Kwestia godzinki może dwóch.
Jeżeli chcesz być na bieżąco, uderz do mnie bezpośrednio sqlek AT sqlek DOT org

## Uruchomienie i testy

  docker compose build
  docker compose up
  go test ./...

### Uwaga na windowsie

Testy zawarte w tym projekcie wymagają dockera na maszynie.
DockerTest obsługuje windowsa, ale jeszcze nie obsługuje TLS

Prawoklik ikona dockera, ustawienia i załączyć
> Expose daemon on tcp://localhost:2375 without TLS

## Rest API

Pełen opis znajdziesz w ZadanieGo.pdf.
Jak starczy czasu opiszę je też tutaj.

### GET /api/urls

Listuje adresy które aktualnie są obserwowane.

### POST /api/urls

Dodaje adres to cyklicznego obserwowania.

### PATCH /api/urls/{ID}

Modyfikuje jeden z obserwowanych adresów.

### DELETE /api/urls/{ID}

Przestaje obserwować dany adres i usuwa wzmianki z bazy danych.

### GET /api/urls/{ID}/history

Listuje historię podglądu danego adresu.

## Odpowiedzi na pytania otwarte

* Jakie problemy mogą się pojawić w trakcie działania serwisu?
* W jaki sposób zabezpieczyć serwis przed złośliwymi zapytaniami?

Ten projekt, nawet perfekcyjnie wykonany ma potencjał na nadużycia,
czy wykorzystaniu przez osoby trzecie do DoS.

Z pomysłów jakie mam na zminimalizowanie powierzchni ataku to:

* Dodanie ograniczenia dla kolumny url "unique" aby jeden adres nie mógł się pojawić wielokrotnie
* Normalizacja adresów aby uniknać nadużyć czy pomyłek w rodzaju `example.com?foo=bar` i `example.com?foo=baz`.
  Ograniczenie
* Whitelista z adresami/domenami które przewidujemy obserwować.

### Ile zapytań i danych serwis jest w stanie obsłużyć?

Ile jednoczesnych obserwacji można by było prowadzić na raz?
Więcej niż jedno. Ile dokładnie zależy od maszyny.

### W jaki sposób można się przygotować na większy ruch?

Mój pomysł na zwiększenie przepustowości systemu opisałem w `Moja opinia co do serwer/worker`.

### W jaki sposób sprawdzić, że wszystko działa? Co monitorować?

(TODO)

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
  Parę dni temu robiłem zadanie z Cassandrą i była by to fajna opcja.
  Można byłoby rozważyć Redis, Kafka czy Influx ale nie mam na nich doświadczenia.

* Ruch dodaj/usuń zadanie zostawiłbym na bazie SQL,
  ponieważ nie sądzę aby te czynności były wykonywane często.
  Nie wiem jak MySQL ale postgresa można shardować w stylu master-slaves dla reduntancji.

## Dług techniczny

Doba nie z gumy, deadline to deadline.
W moich projektach utrzymuję listę
`nie do końca eleganckie, działa, do poprawki jak będzie luźniej`
Oto taka lista dla tego projektu

### GORM.createdAt

Jak embeduje się gorm.Model to delete nie jest twarde i by zostawały w bazie.
Bez osadzenia gorm.Model, wbudowane createdAt przestało działać.
Jestem przekonany że jakaś pierdółka.

ETA jedna dniówka

### Zmigrować testy integracyjne na DockerTest

Trzeba by było się wgryść jak zlinkować poszczególne kontenery,
przy użyciu DockerTest w podobny sposób w jaki aktualnie jest ustawiony
docker-compose.

Przy okazji można by było zrefaktorować testy.
Wyciagnąć main do katalogu cmd, testy umieścić w katalogu głównym.

ETA dniówka - trzy

### DockerTest a TLS na windows

DockerTest bardzo przydatne narzedzie.
Zamiast mockować można testować z prawdziwą bazą i prawdziwym driverem.

Chciałbym obczaić jak używać tego narzędzia z docker compose czy docker swarm.
