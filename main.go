package main

import (
	"fmt"
	"os"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

// come per bootsrap lo spazio massimo delle colonne in maroto è 12
// invece per la riga deve essere solo impostata la lunghezza di essa
// extrapolate se messo a false torna contenuto text a capo automaticamente( inserisci altre righre autoamticamente)

type MovimentoA struct {
	ID                   int
	RegistroID           int
	Origine              string
	CantiereId           *int
	SottocatenaID        *int
	Data                 time.Time
	Tipo                 string
	Numero               int
	CodCer               string
	UnitaMisura          string
	StatoFisico          string
	Pericolosita         []string
	FormularioNumero     *string
	FormularioData       time.Time
	ProduttoreID         *int
	DestinatarioID       *int
	TrasportatoreID      *int
	IntermediarioID      *int
	Quantita             float64
	DataInserimento      time.Time
	Stampato             bool
	Annullato            bool
	Note                 string
	Riferimenti          []MovimentoA
	AttivitaDestinazione *string
}

func main() {
	numeroFormulario := "20202"
	array := [3]MovimentoA{
		{ID: 1, RegistroID: 10, Origine: "Italy", Data: time.Now(), Tipo: "Carico", Numero: 1, CodCer: "2333", UnitaMisura: "kg", StatoFisico: " 2 - Solido non pulverulento ", Pericolosita: []string{"Hp14"}, FormularioNumero: &numeroFormulario, FormularioData: time.Now(), Quantita: 5000, DataInserimento: time.Now(), Note: "FILTRI MOTORE"},
		{ID: 2, RegistroID: 20, Origine: "France", Data: time.Now(), Tipo: "Carico", Numero: 2, CodCer: "5666", UnitaMisura: "kg", StatoFisico: " 2 - Solido non pulverulento", Pericolosita: []string{"Hp12"}, FormularioNumero: &numeroFormulario, FormularioData: time.Now(), Quantita: 6000, DataInserimento: time.Now(), Note: "Quantità residua ultimo carico (12/2021):900 kg"},
		{ID: 3, RegistroID: 30, Origine: "Germany", Data: time.Now(), Tipo: "Scarico", Numero: 3, CodCer: "2333", UnitaMisura: "kg", StatoFisico: "2 - Solido non pulverulento", Pericolosita: []string{"Hp14"}, FormularioNumero: &numeroFormulario, FormularioData: time.Now(), Quantita: 80000, DataInserimento: time.Now(), Note: "FILTRI MOTORE Quantità residua ultimo carico (24/2021): 20 kg"},
	}
	array[2].Riferimenti = append(array[2].Riferimenti, array[0], array[1])
	begin := time.Now()
	m := pdf.NewMaroto(consts.Landscape, consts.A4)
	m.SetBorder(true)
	// tableHeadings := []string{"Carico", "Caratteristiche rifiuto", "Quantita", "Luogo di produzione e attività", "Annotazioni"}
	// contents := [][]string{{"Apple", "Red and juicy", "2.00", "ciao", "prova"}, {"Apple", "Red and juicy", "2.00", "ciao", "prova"}}
	for _, v := range array {

		dataMovimento := v.Data.Format("2006-01-02")
		stringaDataMovimento := fmt.Sprint("Del       ", dataMovimento)
		annoMovimento := v.Data.Year()
		stringaAnno := fmt.Sprint(annoMovimento)
		stringaNumero := fmt.Sprint("N.         ", v.Numero, "/", stringaAnno)
		stringaNumeroFormulario := fmt.Sprint("N.         ", *v.FormularioNumero)
		DataFormulario := v.FormularioData.Format("2006-01-02")
		stringaDataFormulario := fmt.Sprint("Del       ", DataFormulario)
		stringaCer := fmt.Sprint("C.E.R.: ", v.CodCer)
		stringaStatofisico := fmt.Sprint("Stato fisico: ", v.StatoFisico)
		Descrizione := "Descrizione: * assorbenti, materiali filtranti (inclusi filtri dell'olio non specificati altrimenti), stracci e indumenti protettivi, contaminati da sostanze pericolose"
		var stringaPericolosità string
		stringaPericolosità = "Classi di pericolosità: "
		if len(v.Pericolosita) > 0 {
			for _, v := range v.Pericolosita {
				stringaPericolosità = stringaPericolosità + v
			}
		}
		Destinazione := "Destinazione: R04 - Riciclo/recupero dei metalli e dei composti metallici"
		stringaQuantita := fmt.Sprint(v.Quantita)
		m.Row(50, func() {
			// PRIMA COLONNA
			m.Col(2, func() {
				//
				m.Text(v.Tipo, props.Text{
					Top:   3,
					Style: consts.Bold,
					Size:  8,
					Left:  1,
				})
				//
				m.Text(stringaDataMovimento, props.Text{
					Top:  7,
					Size: 8,
					Left: 1,
				})
				//
				m.Text(stringaNumero, props.Text{
					Top:  10,
					Size: 8,
					Left: 1,
				})
				//
				m.Text("Formulario", props.Text{
					Top:   13,
					Style: consts.Bold,
					Size:  8,
					Left:  1,
				})
				//
				m.Text(stringaNumeroFormulario, props.Text{
					Top:  16,
					Size: 8,
					Left: 1,
				})
				//
				m.Text(stringaDataFormulario, props.Text{
					Top:  19,
					Size: 8,
					Left: 1,
				})
				//
				if v.Tipo == "Scarico" {
					m.Text("Rif. op di scarico", props.Text{
						Top:   22,
						Style: consts.Bold,
						Left:  1,
						Size:  8,
					})
					top := 25.0
					for _, v := range v.Riferimenti {
						m.Text(faiStringaNumero(v, stringaAnno), props.Text{
							Top:  top,
							Size: 7,
							Left: 1,
						})
						top += 3
					}
				}
			})
			// COLONNA CARATTERISTICHE RIFIUTO
			m.Col(3, func() {
				//
				m.Text("Caratteristiche rifiuto", props.Text{
					Top:   3,
					Style: consts.Bold,
					Align: consts.Center,
					Size:  8,
				})
				//
				m.Text(stringaCer, props.Text{
					Top:  7,
					Size: 8,
					Left: 1,
				})
				// descrizione se maggiore di un certo tot metto extarpolate false dato che voglio che contenuto che resta fuori me lo amnda a capo
				m.Text(Descrizione, props.Text{
					Top:         10,
					Size:        8,
					Left:        1,
					Extrapolate: false,
				})
				// vale  la stessa cosa di descriozne
				m.Text(stringaStatofisico, props.Text{
					Top:         13,
					Size:        8,
					Left:        1,
					Extrapolate: false,
				})
				// vale  la stessa cosa di descriozne
				m.Text(stringaPericolosità, props.Text{
					Top:         16,
					Size:        8,
					Left:        1,
					Extrapolate: false,
				})
				// vale  la stessa cosa di descriozne
				m.Text(Destinazione, props.Text{
					Top:         19,
					Size:        8,
					Left:        1,
					Extrapolate: false,
				})
			})
			// COLONNA QUANTITA
			m.Col(2, func() {
				m.Text("Quantità", props.Text{
					Top:   3,
					Style: consts.Bold,
					Align: consts.Center,
					Size:  8,
				})
				m.Text(v.UnitaMisura, props.Text{
					Top:   8,
					Size:  8,
					Align: consts.Center,
				})
				m.Text(stringaQuantita, props.Text{
					Top:   11,
					Size:  8,
					Align: consts.Center,
				})
			})
			// COLONNA LUOGO DI PRODUZIONE E ATTIVITà DI PROVENIENZA
			m.Col(3, func() {
				m.Text("Luogo di produzione e attivita di provenienza", props.Text{
					Top:   3,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
				})
				m.Text("Intermediario/commerciante", props.Text{
					Top:   10,
					Style: consts.Bold,
					Size:  8,
					Align: consts.Center,
					Left:  1,
				})
				m.Text("Denominazione", props.Text{
					Top:  13,
					Size: 8,
					Left: 1,
				})
				m.Text("Sede", props.Text{
					Top:  16,
					Size: 8,
					Left: 1,
				})
				m.Text("C.F", props.Text{
					Top:  19,
					Size: 8,
					Left: 1,
				})
				m.Text("Iscrizione Albo", props.Text{
					Top:  22,
					Size: 8,
					Left: 1,
				})

			})
			// COLONNA ANNOTAZIONI
			m.Col(2, func() {
				m.Text("Annotazioni", props.Text{
					Top:   3,
					Style: consts.Bold,
					Align: consts.Center,
					Size:  8,
				})
				if v.Note != "" {
					m.Text(v.Note, props.Text{
						Top:         7,
						Left:        1,
						Size:        8,
						Extrapolate: false,
					})
				}
			})
		})
	}
	// m.TableList(tableHeadings, contents, props.TableList{
	// 	HeaderProp: props.TableListContent{
	// 		Size:      9,
	// 		GridSizes: []uint{2, 3, 1, 4, 2},
	// 	},
	// 	ContentProp: props.TableListContent{
	// 		Size:      8,
	// 		GridSizes: []uint{2, 3, 1, 4, 2},
	// 	},
	// 	Align:                consts.Center,
	// 	AlternatedBackground: &color.Color{Red: 30, Green: 50, Blue: 20},
	// 	Line:                 false,
	// })
	err := m.OutputFileAndClose("pdfs/zpl.pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
		os.Exit(1)
	}

	end := time.Now()
	fmt.Println(end.Sub(begin))
}

func faiStringaNumero(movimento MovimentoA, annoMovimento string) string {
	return fmt.Sprint(movimento.Numero, "/", annoMovimento)
}
