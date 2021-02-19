package main

//Importar librerias
//api de contratos inteligentes nuevo patron de diseño de 
//contratos inteligentes
//"github.com/hyperledger/fabric-contract-api-go/contractapi"

import (
        "encoding/json"
	"fmt"
        "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

//Definicion de estructuras
//Forma de usar la estructura de contratos inteligentes
type SmartContract struct {
        contractapi.Contract
}

//Crear un conjunto de estructuras
//Definir activos Demandante, TipoDocumento y ContenidoDocumento

type Documento struct {
        Demandante         string `json:"demandante"`
        TipoDocumento      string `json:"tipoDocumento"`
        ContenidoDocumento string `json:"contenidoDocumento"`
}

//Funciones

//co s *SmartContract nos aseguramos que la función sea parte de la 
// estructura smartcontract
//ctx: para gestionar la transacción
//las demas variables de tipo string son propiedades de los activos
//Esta función permite almacenar en la red blockchain de pjud.

func (s *SmartContract) Set(ctx contractapi.TransactionContextInterface, documentoId string, demandante string, tipoDocumento string, contenidoDocumento string) error {
    //Validar Existencia de transacción
    doc, err := s.Query(ctx, documentoId)

    if doc != nil{
        fmt.Printf("documentoId ya existe, error: %s", err.Error())
        return err
    }

    //estructura de go
    doc = &Documento{
              Demandante: demandante,
              TipoDocumento: tipoDocumento,
              ContenidoDocumento: contenidoDocumento,
    }
    //transformar a byte al elemento documento
    documentoAsBytes, err := json.Marshal(doc)
    if err != nil {
       fmt.Printf("Marshal error: %s", err.Error())
       return err
    }
    //PutState permite guardar en el libro distribuido el id y valor
    //el libro es de tipo clave valor distribuido
    return ctx.GetStub().PutState(documentoId, documentoAsBytes)

}

//Función para consultar en el ledger o libro mayor.

func (s *SmartContract) Query(ctx contractapi.TransactionContextInterface, documentoId string) (*Documento, error) {
   //Consultar estado actual esta relacionado con un estampa de tiempo y 
   //transacción de id
   //Historial de los activos
   //estampa de tiempo: guardar un estado de tiempo en blockchain 
   //se guarda de forma automática

   //Retorna valor en bytes y el contenido del posible error
   documentoAsBytes, err := ctx.GetStub().GetState(documentoId)

   if err != nil {
      return nil, fmt.Errorf("No se puede leer del estado mundial. %s", err.Error())
   }

   if documentoAsBytes == nil {
      return nil, fmt.Errorf("%s No existe", documentoId)
   }

   doc := new(Documento)
   err = json.Unmarshal(documentoAsBytes, doc)
   if err != nil {
      return nil, fmt.Errorf("Unmarshal error. %s", err.Error())
   }

   return doc, nil
}

func main() {
    //Inicializar contrato inteligente
    chaincode, err := contractapi.NewChaincode(new(SmartContract))
    //Diferente a nulo
    if err != nil {
      fmt.Printf("Error al crear chaincode envio-documento-demanda-escrito: %s", err.Error())
      return
    }

    //Iniciar contrato inteligente y verificar que no tenga errores
    if err := chaincode.Start(); err != nil {
       fmt.Printf("Error starting envio-documento-demanda-escrito chaincode: %s", err.Error())
    }
}
