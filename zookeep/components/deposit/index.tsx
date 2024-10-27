'use client'
import React, { useState } from 'react';
import { Button } from '../ui/button';
import { Input } from '../ui/input';
import { Select, SelectTrigger, SelectValue, SelectContent, SelectItem } from '../ui/select'
import { useTranslation } from '@/app/i18n/client';
import Footer from '../footer';

interface TransProps {
    lng:string
}

function TransactionForm({lng}:TransProps) {
    const [amount, setAmount] = useState('');
    const [transactionType, setTransactionType] = useState('deposit'); // 'deposit' or 'withdraw'

    const {t} = useTranslation(lng,"home",undefined)
    const handleSubmit = (e) => {
        e.preventDefault();
        console.log(`Transaction Type: ${transactionType}, Amount: ${amount}`);
        setAmount('');
    };

    return (
        <>
        <div className="p-4 max-w-md mx-auto">
            <h2 className="text-xl font-bold mb-4">{transactionType === 'deposit' ? t('deposit') : t('withdraw')}</h2>
            <form onSubmit={handleSubmit} className="space-y-4">
                <div>
                    <label htmlFor="amount" className="block text-sm font-medium">{transactionType === 'deposit'?t('deposit'):t('withdraw')}</label>
                    <Input
                        type="number"
                        id="amount"
                        value={amount}
                        onChange={(e) => setAmount(e.target.value)}
                        required
                        className="mt-1"
                    />
                </div>
                <div>
                    <label htmlFor="transactionType" className="block text-sm font-medium">{t('transactionType')}</label>
                    <Select
                        value={transactionType}
                        onValueChange={setTransactionType}
                        className="mt-1"
                    >
                        <SelectTrigger>
                            <SelectValue placeholder="เลือกประเภท" />
                        </SelectTrigger>
                        <SelectContent>
                            <SelectItem value="deposit">ฝากเงิน</SelectItem>
                            <SelectItem value="withdraw">ถอนเงิน</SelectItem>
                        </SelectContent>
                    </Select>
                </div>
                <Button type="submit" className="w-full">
                    {transactionType === 'deposit' ? 'ฝากเงิน' : 'ถอนเงิน'}
                </Button>
            </form>
        </div>
       {/* <Footer lng={lng} />  */}
        </>
    );
};

export default TransactionForm;