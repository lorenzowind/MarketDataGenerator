import csv
import random
from datetime import datetime, timedelta

# Constantes que podem ser modificadas
ASSET_NAME = 'PETR4'  # Nome do ticker do ativo
AVERAGE_SPREAD = 0.05  # Spread médio em R$
TRADING_START_TIME = datetime.strptime('09:00:00', '%H:%M:%S')
TRADING_END_TIME = datetime.strptime('17:00:00', '%H:%M:%S')

# Configurações adicionais
NUM_OFFERS = 1000  # Número total de ofertas para gerar
CANCEL_PROBABILITY = 0.05  # Probabilidade de uma oferta ser cancelada
EDIT_PROBABILITY = 0.1  # Probabilidade de uma oferta ser editada
PRICE_BASE = 20.00  # Preço base do ativo
PRICE_VOLATILITY = 1.00  # Volatilidade do preço
QUANTITY_BASE = 100  # Quantidade base das ofertas
QUANTITY_VARIATION = 50  # Variação na quantidade das ofertas

# Listas para armazenar as ofertas e negócios
buy_offers = []
sell_offers = []
trades = []

# Contadores e IDs
offer_generation_id = 1
offer_sequence_id = 1
trade_id = 1

# Funções auxiliares
def generate_timestamp():
    total_seconds = (TRADING_END_TIME - TRADING_START_TIME).seconds
    random_seconds = random.randint(0, total_seconds)
    timestamp = TRADING_START_TIME + timedelta(seconds=random_seconds)
    return timestamp

def generate_offer(is_buy):
    global offer_generation_id, offer_sequence_id

    # Determinar o preço da oferta
    if is_buy:
        price = PRICE_BASE - AVERAGE_SPREAD/2 + random.uniform(-PRICE_VOLATILITY, PRICE_VOLATILITY)
    else:
        price = PRICE_BASE + AVERAGE_SPREAD/2 + random.uniform(-PRICE_VOLATILITY, PRICE_VOLATILITY)
    price = round(price, 2)

    # Gerar quantidade
    quantity_total = QUANTITY_BASE + random.randint(-QUANTITY_VARIATION, QUANTITY_VARIATION)
    quantity_total = max(1, quantity_total)

    # Gerar oferta
    offer = {
        'cod_evento_oferta': '0',  # Criação
        'cod_simbolo_instrumento_negociacao': ASSET_NAME,
        'dthr_inclusao_oferta': generate_timestamp(),
        'num_geracao_oferta': offer_generation_id,
        'num_identificacao_conta': random.randint(1000, 9999),
        'num_negocio': 0,
        'num_sequencia_oferta': offer_sequence_id,
        'num_sequencia_oferta_secundaria': offer_sequence_id,
        'qte_divulgada_oferta': quantity_total,
        'qte_negociada': 0,
        'qte_total_oferta': quantity_total,
        'val_preco_oferta': price
    }

    offer_generation_id += 1
    offer_sequence_id += 1

    return offer

def match_offers():
    global trade_id, offer_generation_id, offer_sequence_id

    # Continuar casando ofertas enquanto houver correspondência
    while True:
        # Filtrar ofertas ativas (não canceladas e com quantidade disponível)
        active_buy_offers = [o for o in buy_offers if o['cod_evento_oferta'] != '4' and o['qte_divulgada_oferta'] > 0]
        active_sell_offers = [o for o in sell_offers if o['cod_evento_oferta'] != '4' and o['qte_divulgada_oferta'] > 0]

        buy_offers_sorted = sorted(active_buy_offers, key=lambda x: (-x['val_preco_oferta'], x['dthr_inclusao_oferta']))
        sell_offers_sorted = sorted(active_sell_offers, key=lambda x: (x['val_preco_oferta'], x['dthr_inclusao_oferta']))

        trade_occurred = False

        for buy in buy_offers_sorted:
            for sell in sell_offers_sorted:
                if buy['val_preco_oferta'] >= sell['val_preco_oferta']:
                    # Realizar negócio
                    trade_quantity = min(buy['qte_divulgada_oferta'], sell['qte_divulgada_oferta'])
                    trade_price = sell['val_preco_oferta']  # Prioridade para o vendedor

                    # Atualizar ofertas
                    buy['qte_divulgada_oferta'] -= trade_quantity
                    buy['qte_negociada'] += trade_quantity
                    sell['qte_divulgada_oferta'] -= trade_quantity
                    sell['qte_negociada'] += trade_quantity

                    # Atualizar num_negocio
                    buy['num_negocio'] = trade_id
                    sell['num_negocio'] = trade_id

                    # Atualizar cod_evento_oferta se a oferta foi totalmente executada
                    if buy['qte_divulgada_oferta'] == 0:
                        buy['cod_evento_oferta'] = 'F'  # Oferta totalmente negociada
                    else:
                        buy['cod_evento_oferta'] = '0'  # Ainda ativa

                    if sell['qte_divulgada_oferta'] == 0:
                        sell['cod_evento_oferta'] = 'F'
                    else:
                        sell['cod_evento_oferta'] = '0'

                    # Criar registro de negócio para a compra
                    trade_buy = {
                        'cod_natureza_operacao': 'C',
                        'cod_simbolo_instrumento_negociacao': ASSET_NAME,
                        'dthr_negocio': max(buy['dthr_inclusao_oferta'], sell['dthr_inclusao_oferta']).isoformat(),
                        'num_geracao_oferta': buy['num_geracao_oferta'],
                        'num_identificacao_conta': buy['num_identificacao_conta'],
                        'num_negocio': trade_id,
                        'num_sequencia_oferta': buy['num_sequencia_oferta'],
                        'num_sequencia_oferta_secundaria': buy['num_sequencia_oferta_secundaria'],
                        'qte_negocio': trade_quantity,
                        'val_preco_negocio': trade_price
                    }

                    # Criar registro de negócio para a venda
                    trade_sell = {
                        'cod_natureza_operacao': 'V',
                        'cod_simbolo_instrumento_negociacao': ASSET_NAME,
                        'dthr_negocio': max(buy['dthr_inclusao_oferta'], sell['dthr_inclusao_oferta']).isoformat(),
                        'num_geracao_oferta': sell['num_geracao_oferta'],
                        'num_identificacao_conta': sell['num_identificacao_conta'],
                        'num_negocio': trade_id,
                        'num_sequencia_oferta': sell['num_sequencia_oferta'],
                        'num_sequencia_oferta_secundaria': sell['num_sequencia_oferta_secundaria'],
                        'qte_negocio': trade_quantity,
                        'val_preco_negocio': trade_price
                    }

                    trades.append(trade_buy)
                    trades.append(trade_sell)

                    trade_id += 1
                    trade_occurred = True

                    # Atualizar num_geracao_oferta
                    offer_generation_id += 1
                    buy['num_geracao_oferta'] = offer_generation_id
                    sell['num_geracao_oferta'] = offer_generation_id

                    # Se a oferta foi totalmente executada, não precisa continuar tentando casar
                    if buy['qte_divulgada_oferta'] == 0:
                        break  # Passar para a próxima oferta de compra

            # Se houve um negócio, reiniciar o loop para reordenar as ofertas
            if trade_occurred:
                break

        # Se não houve nenhum negócio neste ciclo, interromper o loop
        if not trade_occurred:
            break

# Gera ofertas iniciais
for _ in range(NUM_OFFERS):
    is_buy = random.choice([True, False])
    if is_buy:
        offer = generate_offer(is_buy=True)
        buy_offers.append(offer)
    else:
        offer = generate_offer(is_buy=False)
        sell_offers.append(offer)

    if random.random() < CANCEL_PROBABILITY:
        offer['cod_evento_oferta'] = '4'  # Cancelamento
        offer['qte_divulgada_oferta'] = 0  
    elif random.random() < EDIT_PROBABILITY:
        old_price = offer['val_preco_oferta']
        new_price = old_price + random.uniform(-0.5, 0.5)
        new_price = round(new_price, 2)

        offer_sequence_id += 1
        edited_offer = offer.copy()
        edited_offer['val_preco_oferta'] = new_price
        edited_offer['num_geracao_oferta'] = offer_generation_id
        edited_offer['num_sequencia_oferta_secundaria'] = offer_sequence_id
        edited_offer['dthr_inclusao_oferta'] = generate_timestamp()
        edited_offer['cod_evento_oferta'] = '5'  # Edição
        offer_generation_id += 1

        if is_buy:
            buy_offers.append(edited_offer)
        else:
            sell_offers.append(edited_offer)

match_offers()

def write_csv(filename, data, fieldnames):
    with open(filename, 'w', newline='', encoding='utf-8') as csvfile:
        writer = csv.DictWriter(csvfile, fieldnames=fieldnames, delimiter=';')
        writer.writeheader()
        for row in data:
            if 'dthr_inclusao_oferta' in row:
                row['dthr_inclusao_oferta'] = row['dthr_inclusao_oferta'].isoformat()
            writer.writerow(row)

offer_fieldnames = [
    'cod_evento_oferta',
    'cod_simbolo_instrumento_negociacao',
    'dthr_inclusao_oferta',
    'num_geracao_oferta',
    'num_identificacao_conta',
    'num_negocio',
    'num_sequencia_oferta',
    'num_sequencia_oferta_secundaria',
    'qte_divulgada_oferta',
    'qte_negociada',
    'qte_total_oferta',
    'val_preco_oferta'
]

trade_fieldnames = [
    'cod_natureza_operacao',
    'cod_simbolo_instrumento_negociacao',
    'dthr_negocio',
    'num_geracao_oferta',
    'num_identificacao_conta',
    'num_negocio',
    'num_sequencia_oferta',
    'num_sequencia_oferta_secundaria',
    'qte_negocio',
    'val_preco_negocio'
]

buy_offers_sorted = sorted(buy_offers, key=lambda x: x['dthr_inclusao_oferta'])
sell_offers_sorted = sorted(sell_offers, key=lambda x: x['dthr_inclusao_oferta'])
trades_sorted = sorted(trades, key=lambda x: x['dthr_negocio'])

write_csv('ofertas_compra.csv', buy_offers_sorted, offer_fieldnames)
write_csv('ofertas_venda.csv', sell_offers_sorted, offer_fieldnames)
write_csv('negocios.csv', trades_sorted, trade_fieldnames)

print("Simulação concluída. Arquivos CSV gerados.")
