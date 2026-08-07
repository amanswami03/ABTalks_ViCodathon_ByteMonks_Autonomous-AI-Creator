const variants = {
  primary: 'bg-[#948979] text-[#222831] hover:opacity-90',
  secondary: 'border border-[#948979]/30 bg-transparent text-[#DFD0B8] hover:bg-[#222831]',
  ghost: 'bg-transparent text-[#DFD0B8]/80 hover:bg-[#222831] hover:text-[#DFD0B8]',
};

export default function Button({ children, variant = 'primary', className = '', ...props }) {
  return (
    <button
      className={`rounded-full px-5 py-3 font-semibold transition duration-200 ${variants[variant] || variants.primary} ${className}`.trim()}
      {...props}
    >
      {children}
    </button>
  );
}
