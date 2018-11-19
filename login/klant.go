package main

type Klant struct {
	Knr, Hnr, Ink                    int
	Nm, Vnm, Pc, Hnrt, Gsl, Blg, Rhf string
	Gbd, Krg, Opl, Opm               string
	Brf                              float32
}

type Verkoper struct {
	Mnr  int
	Vanm string
	Vvnm string
}
