import pandas as pd
import random

def simular_pregao_completo(negocios_realizados_df, variacao_preco=0.5, num_ofertas_nao_concretizadas=5):

    ofertas_compra = []
    ofertas_venda = []

    for index, negocio in negocios_realizados_df.iterrows():

        preco_negocio = float(negocio['PrecoNegocio'].replace(',', '.'))
        quantidade = negocio['QuantidadeNegociada']
        hora_fechamento = negocio['HoraFechamento']
        codigo_negocio = negocio['CodigoIdentificadorNegocio']
        codigo_participante_comprador = negocio['CodigoParticipanteComprador']
        codigo_participante_vendedor = negocio['CodigoParticipanteVendedor']
        
        ofertas_compra.append({
            'DataReferencia': negocio['DataReferencia'],
            'CodigoInstrumento': negocio['CodigoInstrumento'],
            'PrecoOferta': preco_negocio, 
            'Quantidade': quantidade,
            'Hora': hora_fechamento,
            'CodigoParticipante': codigo_participante_comprador,
            'TipoOferta': 'Compra',
            'CodigoNegocio': codigo_negocio
        })
        
        ofertas_venda.append({
            'DataReferencia': negocio['DataReferencia'],
            'CodigoInstrumento': negocio['CodigoInstrumento'],
            'PrecoOferta': preco_negocio,
            'Quantidade': quantidade,
            'Hora': hora_fechamento,
            'CodigoParticipante': codigo_participante_vendedor,
            'TipoOferta': 'Venda',
            'CodigoNegocio': codigo_negocio
        })

        for _ in range(num_ofertas_nao_concretizadas):

            ofertas_compra.append({
                'DataReferencia': negocio['DataReferencia'],
                'CodigoInstrumento': negocio['CodigoInstrumento'],
                'PrecoOferta': round(preco_negocio - random.uniform(0.01, variacao_preco), 2),
                'Quantidade': random.choice([100, 200, 300, 400, 500]),
                'Hora': hora_fechamento,
                'CodigoParticipante': random.randint(100, 999),
                'TipoOferta': 'Compra',
                'CodigoNegocio': None
            })
            
            ofertas_venda.append({
                'DataReferencia': negocio['DataReferencia'],
                'CodigoInstrumento': negocio['CodigoInstrumento'],
                'PrecoOferta': round(preco_negocio + random.uniform(0.01, variacao_preco), 2),
                'Quantidade': random.choice([100, 200, 300, 400, 500]),
                'Hora': hora_fechamento,
                'CodigoParticipante': random.randint(100, 999),
                'TipoOferta': 'Venda',
                'CodigoNegocio': None
            })

    df_ofertas_compra = pd.DataFrame(ofertas_compra)
    df_ofertas_venda = pd.DataFrame(ofertas_venda)

    return df_ofertas_compra, df_ofertas_venda

negocios_realizados_df = pd.read_csv('negociosptr4.csv', sep=';')

df_ofertas_compra_completo, df_ofertas_venda_completo = simular_pregao_completo(negocios_realizados_df)

df_ofertas_compra_completo.to_csv('ofertas_compra_geradas.csv', index=False)
df_ofertas_venda_completo.to_csv('ofertas_venda_geradas.csv', index=False)
