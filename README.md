# wp-interview

Zadanie rekrutacyjne dla Wirtualnej Polski.

Zadanie jest aktualnie dopinane, weryfikowane i wysyłane.
Jest udostępnione w takiej formie aby zrobić uprawnienia na githubie.

Kwestia godzinki może dwóch.
Jeżeli chcesz być na bieżąco, uderz do mnie bezpośrednio sqlek AT sqlek DOT org

## Uwaga na windowsie

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

### Moje opinie co do API

Trochę dziwi mnie że nie ma paginacji.
Szczególnie w history by się przydała.

Zastanawia mnie też używanie floatów w jako time.Time i time.Duration.
Realny projekt nie działa w próżni, więc jak jest specyfikacja to się ją trzyma.

## Odpowiedzi na pytania otwarte

## Dług techniczny

Doba nie z gumy, deadline to deadline.
W moich projektach utrzymuję listę
"do poprawki jak będzie luźniej"
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
